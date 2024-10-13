package middlewares

import (
	"log"

	"github.com/danielgtaylor/huma/v2"
)

// RequestLogMiddleware logs the HTTP request method and URL path.
func RequestLogMiddleware(ctx huma.Context, next func(huma.Context)) {
	var partialURL string

	partialURL = ctx.URL().Path

	if ctx.URL().RawQuery != "" {
		partialURL += "?" + ctx.URL().RawQuery
	}

	if ctx.URL().Fragment != "" {
		partialURL += "#" + ctx.URL().Fragment
	}

	log.Println(ctx.Method(), partialURL)

	next(ctx)
}
