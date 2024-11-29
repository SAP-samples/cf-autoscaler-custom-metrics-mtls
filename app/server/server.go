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

	metrics := map[string]interface{}{
		"cpu": &CPUWaster{},
	}

	ah := NewAppHandler(appConfig, "https://autoscaler-metrics-mtls.cf.sap.hana.ondemand.com", metrics)

	router.GET("/", ah.GetHome)

	// custom metrics handlers
	router.GET("/busy/:metricValue", ah.Busy)
	router.GET("/not-busy/:metricValue", ah.NotBusy)

	// cpu handlers
	router.GET("/cpu/:utilization/:minutes", ah.IncreaseCPU)
	router.GET("/cpu/stop", ah.StopCPU)

	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
