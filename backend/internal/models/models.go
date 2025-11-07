package models

// Common response structures
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

type PaginationQuery struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=10"`
}

func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
