package main

import (
	"context"
	"fmt"
	"github.com/aasumitro/pokewar/configs"
	"github.com/aasumitro/pokewar/constants"
	"github.com/aasumitro/pokewar/docs"
	"github.com/aasumitro/pokewar/internal"
	"github.com/aasumitro/pokewar/internal/delivery/middleware"
	"github.com/aasumitro/pokewar/web"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"io"
	"log"
	"net/http"
	"os"
)

// @contact.name 	@aasumitro
// @contact.url 	https://aasumitro.id/
// @contact.email 	hello@aasumitro.id
// @license.name  	MIT
// @license.url   	https://github.com/aasumitro/pokewar/blob/main/LICENSE

var (
	appEngine *gin.Engine
	ctx       = context.Background()
)

func init() {
	configs.LoadEnv()

	configs.Instance.InitDbConn()

	if !configs.Instance.AppDebug {
		gin.SetMode(gin.ReleaseMode)

		accessLogFile, _ := os.Create("./temps/access.log")
		gin.DefaultWriter = io.MultiWriter(accessLogFile, os.Stdout)

		errorLogFile, _ := os.Create("./temps/errors.log")
		gin.DefaultErrorWriter = io.MultiWriter(errorLogFile, os.Stdout)
	}

	appEngine = gin.Default()

	docs.SwaggerInfo.BasePath = appEngine.BasePath()
	docs.SwaggerInfo.Title = configs.Instance.AppName
	docs.SwaggerInfo.Description = "Pocket Monster Battleroyale API Spec"
	docs.SwaggerInfo.Version = configs.Instance.AppVersion
	docs.SwaggerInfo.Host = configs.Instance.AppURL
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func main() {
	appEngine.GET("/", func(context *gin.Context) {
		context.Redirect(http.StatusPermanentRedirect, "/home")
	})

	appEngine.StaticFS("/home",
		http.FS(web.Resource))

	appEngine.GET("/docs/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.DefaultModelsExpandDepth(constants.GinModelsDepth)))

	internal.NewAPIProvider(ctx, appEngine)

	if configs.Instance.AppDebug {
		middleware.RegisterPPROF(appEngine)
	}

	port := os.Getenv("HTTP_PLATFORM_PORT")
	getUrl := func() string {
		if port == "" {
			return configs.Instance.AppURL
		}
		return fmt.Sprintf("127.0.0.1:%s", port)
	}

	log.Fatal(appEngine.Run(getUrl()))
}
