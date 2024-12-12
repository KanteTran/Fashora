package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mvrilo/go-redoc"
	ginredoc "github.com/mvrilo/go-redoc/gin"
)

func SetupApiDocsMiddleware() gin.HandlerFunc {
	doc := redoc.Redoc{
		Title:    "Fashora API Docs",
		SpecFile: "./api-docs/openapi.yaml",
		SpecPath: "/openapi.yaml",
		DocsPath: "/api-docs",
	}

	return ginredoc.New(doc)
}
