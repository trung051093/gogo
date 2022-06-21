package common

type Response struct {
	Data interface{} `json:"data"`
}

type ResponsePagination struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
	Filter     interface{} `json:"filter"`
}

func SuccessResponse(data interface{}) *Response {
	return &Response{Data: data}
}

func SuccessResponsePagination(data interface{}, paging *Pagination, filter interface{}) *ResponsePagination {
	return &ResponsePagination{Data: data, Pagination: paging, Filter: filter}
}
