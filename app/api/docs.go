package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	d "github.com/moonwalker/moonbase/docs"
)

func docs() gin.HandlerFunc {
	d.SwaggerInfo.Version = "1.0"
	d.SwaggerInfo.Host = "moonbase.mw.zone"
	return ginSwagger.WrapHandler(swaggerFiles.Handler)
}
