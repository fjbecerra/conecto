package rest

type Cursor struct {
	Value string
}

type Page[T any] struct {
	Data       []T
	NextCursor *Cursor
	HasMore    bool
}




