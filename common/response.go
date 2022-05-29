package common

type successResponse struct {
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination,omitempty"`
	Filter     interface{} `json:"filter,omitempty"`
}

func SuccessResponse(data, paging, filter interface{}) *successResponse {
	return &successResponse{Data: CompactJson(data), Pagination: paging, Filter: filter}
}

func SimpleSuccessResponse(data interface{}) *successResponse {
	return &successResponse{Data: CompactJson(data), Pagination: nil, Filter: nil}
}
