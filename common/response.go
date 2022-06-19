package common

type Response struct {
	Data interface{} `json:"data"`
}

type ResponsePagination struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Filter     interface{} `json:"filter,omitempty"`
}

func SuccessResponse(data interface{}) *Response {
	return &Response{Data: CompactJson(data)}
}

func SuccessResponsePagination(data interface{}, paging *Pagination, filter interface{}) *ResponsePagination {
	return &ResponsePagination{Data: CompactJson(data), Pagination: paging, Filter: filter}
}
