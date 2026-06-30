package http

type CursorInjector interface {
    Inject(cursor *PageCursor) error
}