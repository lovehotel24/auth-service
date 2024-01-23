package web

type Middleware func(Handler) Handler

// warpMiddleware creates a new handler by wrapping middleware around a final
// handler. The middlewares Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {

	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Lopping backwards ensures that the
	// first middleware of slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
