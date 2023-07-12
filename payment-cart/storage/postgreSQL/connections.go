package postgreSQL

import (
	"cloud.google.com/go/cloudsqlconn"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4/stdlib"
	"hash/crc32"
	"log"
	"net"
	"os"
	"strings"
)

type Connection struct {
	Database       string `json:"database"`
	UseID          string `json:"userID"`
	Password       string `json:"password"`
	ConnectionName string `json:"connectionName"`
	DBIAMUser      string `json:"dbIAMUser"`
}

func NewConnection(ctx context.Context) (*sql.DB, error) {
	var db *sql.DB

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	postgresSecretsName := os.Getenv("CLOUD_SQL_POSTGRES_SECRETS_NAME")
	sqlSecretsReq := fmt.Sprintf("projects/%v/secrets/%v/versions/latest", projectID, postgresSecretsName)
	//log.Printf("sqlSecretsReq: %v\n", sqlSecretsReq)

	payload, err := accessSecretVersion(sqlSecretsReq)

	if err != nil {
		return nil, err
	}

	conn := Connection{}
	if err := json.Unmarshal(payload, &conn); err != nil {
		return nil, fmt.Errorf("error parsing secrets payload, %v", err)
	}

	db, err = connectWithConnector(conn)
	if err != nil {
		return nil, fmt.Errorf("connectIAM: unable to connect: %s", err)
	}
	return db, nil

}

func connectWithConnector(connDetails Connection) (*sql.DB, error) {

	usePrivate, err := getEnv("PRIVATE_IP", false, false)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	dsn := fmt.Sprintf("user=%s password=%s database=%s", connDetails.UseID, connDetails.Password, connDetails.Database)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	var opts []cloudsqlconn.Option
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, err
	}
	// Use the Cloud SQL connector to handle connecting to the instance.
	// This approach does *NOT* require the Cloud SQL proxy.
	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, connDetails.ConnectionName)
	}
	dbURI := stdlib.RegisterConnConfig(config)

	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	return dbPool, nil
}

func accessSecretVersion(name string) ([]byte, error) {
	// name := "projects/my-project/secrets/my-secret/versions/latest"

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create secretmanager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to access secret version: %w", err)
	}

	// Verify the data checksum.
	crc32c := crc32.MakeTable(crc32.Castagnoli)
	checksum := int64(crc32.Checksum(result.Payload.Data, crc32c))
	if checksum != *result.Payload.DataCrc32C {
		return nil, errors.New("data corruption detected")
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	//log.Printf("Plaintext: %s\n", string(result.Payload.Data))

	return result.Payload.Data, nil
}

func getEnv(env string, required bool, mask bool) (string, error) {

	value := os.Getenv(env)
	if len(value) == 0 && required {
		return value, fmt.Errorf("%s environment variable not set", env)
	}
	if mask {
		value = strings.Repeat("*", len(value))
	}
	log.Printf("%v : %v\n", env, value)
	return value, nil

}
