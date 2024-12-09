package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mvrilo/go-redoc"
	ginredoc "github.com/mvrilo/go-redoc/gin"
)

func SetupApiDocs(r *gin.Engine) {
	doc := redoc.Redoc{
		Title:    "Fashora API Docs",
		SpecFile: "./api-docs/openapi.yaml",
		SpecPath: "/openapi.yaml",
		DocsPath: "/api-docs",
	}

	r.GET(doc.SpecPath, ginredoc.New(doc))
	r.GET(doc.DocsPath, ginredoc.New(doc))
}
