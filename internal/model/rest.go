package model

type (
	Response struct {
		Data interface{} `json:"data,omitempty"`
		Meta Meta        `json:"meta,omitempty"`
	}

	Meta struct {
		Code    string `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	}
)

type (
	BadPayloadResponse struct {
		FailedField string
		Tag         string
		Value       string
	}

	PaginationResponse struct {
		Number        int    `json:"number,omitempty"`
		Size          int    `json:"size,omitempty"`
		TotalElements int    `json:"total_elements,omitempty"`
		TotalPages    int    `json:"total_pages,omitempty"`
		Sort          string `json:"sort,omitempty"`
		SortBy        string `json:"sort_by,omitempty"`
	}

	TransactionResponse struct {
		TransactionTimestamp uint   `json:"transaction_timestamp"`
		TransactionId        string `json:"transaction_id"`
	}
)

type (
	SessionRequest struct {
		Id       string `swaggerignore:"true"`
		Partner  string `swaggerignore:"true"`
		Username string `swaggerignore:"true"`
		Msisdn   string `swaggerignore:"true"`
		Email    string `swaggerignore:"true"`
		Role     string `swaggerignore:"true"`
		Fullname string `swaggerignore:"true"`
		ContextRequest
	}

	ContextRequest struct {
		Channel       string `json:"channel" swaggerignore:"true"`
		OS            string `json:"os" swaggerignore:"true"`
		Version       string `json:"version" swaggerignore:"true"`
		Language      string `json:"language" swaggerignore:"true"`
		DeviceId      string `json:"deviceId" swaggerignore:"true"`
		Authorization string `json:"authorization" swaggerignore:"true"`
		RefreshToken  string `json:"refreshToken" swaggerignore:"true"`
		TransactionId string `json:"transactionId" swaggerignore:"true"`
		ApiKey        string `json:"apiKey" swaggerignore:"true"`
	}
)
