package schema

type ProductRequest struct {
	Action   string    `json:"action"`
	Products []Product `json:"products"`
}

type ProductResponse struct {
	ProductRequest
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}
