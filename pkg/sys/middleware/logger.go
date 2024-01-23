package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/lovehotel24/auth-service/pkg/foundation/web"
)

func Logger(log *zap.SugaredLogger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing from context")
			}

			log.Infow("request started", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remote_addr", r.RemoteAddr)

			// Call the next handler.
			err = handler(ctx, w, r)

			log.Infow("request completed", "trace_id", v.TraceID, "method", r.Method, "path", r.URL.Path,
				"remote_addr", r.RemoteAddr, "status_code", v.StatusCode, "since", time.Since(v.Now))

			// Return the error, so it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
