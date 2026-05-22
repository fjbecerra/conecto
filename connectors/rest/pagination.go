package rest

type PageCursor struct {
	Value string
}

type Page[T any] struct {
	Data       []T
	NextCursor *PageCursor
	HasMore    bool
}




