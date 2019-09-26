package econtext

import (
	"context"
	"github.com/labstack/echo"
)

type ContextValue struct {
	TraceId string
}

const ReqContextKey = "reqcontextkey"

func ParseContext(c echo.Context) context.Context {
	ctx := context.WithValue(context.Background(), ReqContextKey, ContextValue{c.Request().Header.Get("X-Micro-Trace-Id")})
	return ctx
}
