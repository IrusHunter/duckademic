package platform

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/IrusHunter/duckademic/shared/contextutil"
	"github.com/IrusHunter/duckademic/shared/jsonutil"
	"github.com/IrusHunter/duckademic/shared/logger"
)

// HandlerFunc defines the signature for HTTP handlers with a context.
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

// Middleware defines the signature for middleware functions that wrap HandlerFunc.
type Middleware func(HandlerFunc) HandlerFunc

// RESTAPIHelper provides utilities for routing, middleware, and logging HTTP requests.
type RESTAPIHelper struct {
	Logger    logger.Logger
	JWTSecret []byte
}

// NewRESTAPIHelper creates a new RESTAPIHelper instance.
//
// It requires a name of the parent class (cn).
func NewRESTAPIHelper(cn string) RESTAPIHelper {
	return RESTAPIHelper{
		Logger: logger.NewLogger(cn+".txt", cn),
	}
}

func NewRESTAPIHelperWithAuth(cn string, jwtSecret []byte) RESTAPIHelper {
	return RESTAPIHelper{
		Logger:    logger.NewLogger(cn+".txt", cn),
		JWTSecret: jwtSecret,
	}
}

// NewHandler chains the given middlewares around a handler in the specified order.
func (rh *RESTAPIHelper) NewHandler(h HandlerFunc, m ...Middleware) HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

// TraceMiddleware sets a trace ID in the context for request tracing.
func (rh *RESTAPIHelper) TraceMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		ctx = contextutil.SetTraceID(r.Context())
		next(ctx, w, r)
	}
}

// LoggingMiddleware logs the incoming request and the response status, size, and duration.
func (rh *RESTAPIHelper) LoggingMiddleware(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		traceID := contextutil.GetTraceID(ctx)

		rh.Logger.Log(traceID, "LoggingMiddleware",
			fmt.Sprintf("request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr),
			logger.MiddlewareRequestReceived,
		)

		start := time.Now()

		rw := &responseWriter{ResponseWriter: w}

		next(ctx, rw, r)

		rh.Logger.Log(traceID, "LoggingMiddleware",
			fmt.Sprintf("response %d %s %db %s", rw.status, strings.ToUpper(http.StatusText(rw.status)), rw.size, time.Since(start)),
			logger.MiddlewareRequestFinished,
		)
	}
}

func (rh *RESTAPIHelper) NewAuthMiddleware(requiredPermissions []string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			traceID := contextutil.GetTraceID(ctx)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				jsonutil.ResponseWithError(
					w,
					http.StatusUnauthorized,
					rh.Logger.LogAndReturnError(
						traceID,
						"AuthMiddleware",
						fmt.Errorf("missing authorization header"),
						logger.MiddlewareFailed,
					),
				)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				jsonutil.ResponseWithError(
					w,
					http.StatusUnauthorized,
					rh.Logger.LogAndReturnError(
						traceID,
						"AuthMiddleware",
						fmt.Errorf("invalid authorization header format"),
						logger.MiddlewareFailed,
					),
				)
				return
			}

			token := parts[1]

			claims, err := jsonutil.ParseAccessTokenWithValidation(token, rh.JWTSecret)
			if err != nil {
				jsonutil.ResponseWithError(
					w,
					http.StatusUnauthorized,
					rh.Logger.LogAndReturnError(
						traceID,
						"AuthMiddleware",
						fmt.Errorf("invalid token: %w", err),
						logger.MiddlewareFailed,
					),
				)
				return
			}

			if len(requiredPermissions) > 0 {
				if !hasPermissions(claims.Permissions, requiredPermissions) {
					jsonutil.ResponseWithError(
						w,
						http.StatusForbidden,
						rh.Logger.LogAndReturnError(
							traceID,
							"AuthMiddleware",
							fmt.Errorf("missing required permissions"),
							logger.MiddlewareFailed,
						),
					)
					return
				}
			}

			ctx = contextutil.SetAccessClaims(ctx, claims)
			next(ctx, w, r)
		}
	}
}

func hasPermissions(userPerms []string, required []string) bool {
	if len(required) == 0 {
		return true
	}

	permMap := make(map[string]struct{}, len(userPerms))
	for _, p := range userPerms {
		permMap[p] = struct{}{}
	}

	for _, req := range required {
		if _, ok := permMap[req]; !ok {
			return false
		}
	}

	return true
}

// NewDefaultHandler creates a handler wrapped with TraceMiddleware and LoggingMiddleware.
func (rh *RESTAPIHelper) NewDefaultHandler(h HandlerFunc) HandlerFunc {
	return rh.NewHandler(h, rh.TraceMiddleware, rh.LoggingMiddleware)
}

func (rh *RESTAPIHelper) NewDefaultHandlerWithAuth(h HandlerFunc, requiredPermissions []string) HandlerFunc {
	return rh.NewHandler(h, rh.TraceMiddleware, rh.LoggingMiddleware, rh.NewAuthMiddleware(requiredPermissions))
}

// NewRoute registers a set of HTTP methods and their corresponding handlers for a given path.
// If the method is not allowed, it responds with a 405 status and sets the "Allow" header.
func (rh *RESTAPIHelper) NewRoute(path string, routes map[string]HandlerFunc) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		f, ok := routes[r.Method]
		if !ok {
			keys := rh.getAllowedMethods(routes)
			w.Header().Set("Allow", keys)
			jsonutil.ResponseWithError(w, 405, rh.Logger.LogAndReturnError("", "Route",
				fmt.Errorf("method %q not allowed (allowed methods %s)", r.Method, rh.getAllowedMethods(routes)),
				logger.HandlerBadRequest))
			return
		}

		f(r.Context(), w, r)
	})
}

func (rh *RESTAPIHelper) getAllowedMethods(m map[string]HandlerFunc) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return strings.Join(keys, ", ")
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}
