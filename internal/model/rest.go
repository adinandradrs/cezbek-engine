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
		TransactionTimestamp int64  `json:"transaction_timestamp"`
		TransactionId        string `json:"transaction_id"`
	}

	SessionResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
		Expired      *int64 `json:"expired,omitempty"`
	}
)

type (
	FindByIdRequest struct {
		Id int64 `json:"id"`
		SessionRequest
	}

	SearchRequest struct {
		TextSearch string `json:"text_search"`
		Start      uint   `json:"start" binding:"required" example:"0"`
		Limit      uint   `json:"limit" binding:"required" example:"5"`
		SortBy     string `json:"sort_by" `
		Sort       string `json:"sort" enums:"ASC,DESC"`
		SessionRequest
	}

	SessionRequest struct {
		Id          int64  `swaggerignore:"true"`
		PartnerCode string `swaggerignore:"true"`
		Username    string `swaggerignore:"true"`
		Msisdn      string `swaggerignore:"true"`
		Email       string `swaggerignore:"true"`
		Role        string `swaggerignore:"true"`
		Fullname    string `swaggerignore:"true"`
		ContextRequest
	}

	ContextRequest struct {
		Channel       string `json:"channel,omitempty" swaggerignore:"true"`
		OS            string `json:"os,omitempty" swaggerignore:"true"`
		Version       string `json:"version,omitempty" swaggerignore:"true"`
		DeviceId      string `json:"device_id,omitempty" swaggerignore:"true"`
		Authorization string `json:"authorization,omitempty" swaggerignore:"true"`
		AuthSignature string `json:"auth_signature,omitempty" swaggerignore:"true"`
		RefreshToken  string `json:"refresh_token,omitempty" swaggerignore:"true"`
		TransactionId string `json:"transaction_id,omitempty" swaggerignore:"true"`
		ApiKey        string `json:"api_key,omitempty" swaggerignore:"true"`
	}
)
