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
		Number        int    `json:"number,omitempty" example:"1"`
		Size          int    `json:"size,omitempty" example:"10"`
		TotalElements int    `json:"total_elements,omitempty" example:"100"`
		TotalPages    int    `json:"total_pages,omitempty" example:"10"`
		Sort          string `json:"sort,omitempty" example:"ASC"`
		SortBy        string `json:"sort_by,omitempty" example:"id"`
	}

	TransactionResponse struct {
		TransactionTimestamp int64  `json:"transaction_timestamp" example:"11285736234"`
		TransactionId        string `json:"transaction_id" example:"TRX0012345678"`
	}

	SessionResponse struct {
		Token        string `json:"token,omitempty" example:"**secret**"`
		RefreshToken string `json:"refresh_token,omitempty" example:"**secret**"`
		AccessToken  string `json:"access_token,omitempty" example:"**secret**"`
		Expired      *int64 `json:"expired,omitempty" example:"11234823643"`
	}
)

type (
	FindByIdRequest struct {
		Id int64 `json:"id"`
		SessionRequest
	}

	SearchRequest struct {
		TextSearch string `json:"text_search"`
		Start      int    `json:"start" binding:"required" example:"0"`
		Limit      int    `json:"limit" binding:"required" example:"5"`
		SortBy     string `json:"sort_by" `
		Sort       string `json:"sort" enums:"ASC,DESC"`
		SessionRequest
	}

	SessionRequest struct {
		Id       int64  `swaggerignore:"true"`
		Username string `swaggerignore:"true"`
		Msisdn   string `swaggerignore:"true"`
		Email    string `swaggerignore:"true"`
		Role     string `swaggerignore:"true"`
		Fullname string `swaggerignore:"true"`
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

func Page(inp *SearchRequest) {
	if inp.Limit == 0 || inp.Limit > 100 {
		inp.Limit = 10
	}

	if inp.Start <= 1 {
		inp.Start = 0
	} else {
		inp.Start = (inp.Start - 1) * inp.Limit
	}
}

func Pagination(count int, limit int, start int) PaginationResponse {
	return PaginationResponse{
		TotalPages:    (count-1)/int(limit) + 1,
		TotalElements: count,
		Size:          int(limit),
		Number:        int(start),
	}
}
