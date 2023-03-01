package restapi

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	ginlogrus "github.com/toorop/gin-logrus"
)

// SetupRouter setup GIN router
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	//outer := gin.Default()
	router := gin.New()

	router.Use(ginlogrus.Logger(log.StandardLogger()), gin.Recovery())

	c := NewController()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/logs", checkAuthorization(ROLE_READ, ROLE_READWRITE, ROLE_WRITE), c.GetLogs)
		v1.GET("/logs/count", checkAuthorization(ROLE_READ, ROLE_READWRITE, ROLE_WRITE), c.GetLogsCount)
		v1.POST("/logs", checkAuthorization(ROLE_READWRITE, ROLE_WRITE), c.AddLogs)
	}

	return router
}
