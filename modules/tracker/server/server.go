package server

import (
	"fmt"
	"gogo/common"
	"gogo/components/appctx"
	"gogo/components/hasher"
	"gogo/components/mailer"
	"gogo/modules/tracker/middleware"
	"net/http"

	cacheprovider "gogo/components/cache"
	dbprovider "gogo/components/dbprovider"
	esprovider "gogo/components/elasticsearch"
	jaegerprovider "gogo/components/jaeger"
	logprovider "gogo/components/log"
	redisprovider "gogo/components/redis"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

type TrackerServer struct {
	appCtx    appctx.AppContext
	ginEngine *gin.Engine
}

func New() *TrackerServer {
	config := appctx.GetConfig()
	validate := validator.New()

	// graylog
	logprovider.Integrate(config.GetGraylogConfig())

	dbprovider, err := dbprovider.NewDBProvider(
		&dbprovider.DBConfig{
			Host:     config.Database.Host,
			Username: config.Database.Username,
			Password: config.Database.Password,
			Name:     config.Database.Name,
			Port:     config.Database.Port,
			SSLMode:  config.Database.SSLMode,
			TimeZone: config.Database.TimeZone,
		},
		// dbprovider.WithDebug,
		// dbprovider.WithAutoMigration(
		// 	&entity.User{},
		// 	&entity.Auth{},
		// 	&entity.Game{},
		// 	&entity.GameHistory{},
		// ),
	)
	common.PanicIf(err != nil, err)

	configEs := config.GetElasticSearchConfig()
	esService, err := esprovider.NewEsService(configEs)
	common.PanicIf(err != nil, err)

	configRedis := config.GetRedisConfig()
	redisService := redisprovider.NewRedisService(configRedis)
	jaegerService := jaegerprovider.NewExporter(config.GetJaegerConfig())
	cacheService := cacheprovider.NewCacheService(redisService.GetClient())
	mailService := mailer.NewMailer(config.GetMailConfig())
	hashService := hasher.NewHashService()

	appCtx := appctx.NewAppContext(
		dbprovider.GetDBConnection(),
		validate,
		config,
		esService,
		nil,
		redisService,
		nil,
		nil,
		jaegerService,
		cacheService,
		mailService,
		hashService,
	)

	ginEngine := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	ginEngine.Use(cors.New(corsConfig))
	ginEngine.Use(middleware.ErrorHandler(appCtx))
	ginEngine.Use(middleware.SetAppContextIntoRequest(appCtx))

	return &TrackerServer{
		appCtx:    appCtx,
		ginEngine: ginEngine,
	}
}

func (s *TrackerServer) Start() {
	// create routers
	s.createMainRoutes()

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(s.appCtx.GetJaegerService().GetExporter())
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	http.ListenAndServe(fmt.Sprintf(":%d", s.appCtx.GetConfig().Server.Port), &ochttp.Handler{Handler: s.ginEngine})
}
