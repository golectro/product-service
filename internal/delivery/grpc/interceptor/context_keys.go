package interceptor

type ContextKey string

const (
	RequestIDKey   ContextKey = "requestId"
	UserContextKey ContextKey = "auth"
)
