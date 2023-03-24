package server

import (
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer(appConfig *cfenv.App) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	server := &Server{router: router}

	ah := NewAppHandler(appConfig, "https://autoscaler-metrics-mtls.cf.sap.hana.ondemand.com")
	router.GET("/", ah.GetHome)

	// send POST request to custom metrics URL
	router.GET("/busy/:metricValue", ah.Busy)

	router.GET("/not-busy/:metricValue", ah.NotBusy)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
