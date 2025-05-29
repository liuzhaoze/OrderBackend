package ports

const (
	HttpSuccess int64 = iota
	HttpUnknownError
)

const (
	HttpBindRequestBodyError int64 = 1000 + iota
	HttpValidateRequestError
	HttpCreateOrderError
	HttpGetOrderError
)
