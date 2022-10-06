package api

// Pagination contains pagination request
type PaginationRequest struct {
	Offset int32 `form:"offset" binding:"min=0"`
	Limit  int32 `form:"limit,default=10" binding:"min=1,max=100"`
}

// RequestWithID is used for requests searching for a resource by ID
type RequestWithID struct {
	ID int64 `uri:"id" binding:"required,min=1,max=9223372036854775807"`
}
