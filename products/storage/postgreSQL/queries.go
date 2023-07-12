package postgreSQL

const (
	GetProductList = `select product_id, name, price, currency, description, image_url, createtimestamp, updatetimestamp
from products`
)
