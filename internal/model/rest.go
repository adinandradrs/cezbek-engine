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
)
