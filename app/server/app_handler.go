package server

import (
	"bytes"
	"code.cloudfoundry.org/lager"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	cpuLowerBounds = 0
	cpuUpperBounds = 100
)

type AppHandler struct {
	logger               lager.Logger
	appConfig            *cfenv.App
	metrics              map[string]interface{}
	metricsServerBaseUrl string
}

func NewAppHandler(appConfig *cfenv.App, metricsUrl string, metrics map[string]interface{}) *AppHandler {
	return &AppHandler{
		logger:               lager.NewLogger("appHandler"),
		appConfig:            appConfig,
		metrics:              metrics,
		metricsServerBaseUrl: metricsUrl,
	}
}

func (ah *AppHandler) GetHome(context *gin.Context) {
	fmt.Printf("I am GetHome")
	context.JSON(http.StatusOK, gin.H{
		"CF_INSTANCE_KEY":   os.Getenv("CF_INSTANCE_KEY"),
		"CF_INSTANCE_CERT":  os.Getenv("CF_INSTANCE_CERT"),
		"appConfigurations": ah.appConfig,
	})
}

var SubmitScaleUpEventToAutoscaler = func(logger lager.Logger, appConfig *cfenv.App, metricsValue float64, autoscalerApiServerUrl string) (*http.Response, error) {
	return postScaleUpEventToAutoscaler(logger, appConfig, metricsValue, autoscalerApiServerUrl)
}

func (ah *AppHandler) NotBusy(context *gin.Context) {

	ah.sendMetrics("I am not busy with value", context)
}

func (ah *AppHandler) Busy(context *gin.Context) {
	log.Printf("I am busy with value %s", context.Params.ByName("metricValue"))
	ah.sendMetrics("I am busy with value", context)
}

func (ah *AppHandler) IncreaseCPU(context *gin.Context) {

	utilizationParam := context.Param("utilization")
	minutesParam := context.Param("minutes")
	ah.logger.Info("Increasing CPU utilization", lager.Data{"utilization": utilizationParam, "minutes": minutesParam})

	utilization, err := strconv.Atoi(utilizationParam)
	if err != nil || utilization < cpuLowerBounds || utilization > cpuUpperBounds {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Invalid utilization value. It must be between %d and %d.", cpuUpperBounds, cpuUpperBounds),
		})
		return
	}

	minutes, err := strconv.Atoi(minutesParam)
	if err != nil || minutes <= 0 {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid minutes value. It must be a positive integer.",
		})
		return
	}

	cpuWaster, ok := ah.metrics["cpu"].(*CPUWaster)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "CPU waster not configured properly.",
		})
		return
	}
	go func() {
		err := cpuWaster.IncreaseCPU(int64(utilization), time.Duration(minutes)*time.Minute)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Error: %s", err.Error()),
			})
			return
		}
	}()
	if !cpuWaster.isRunning {
		context.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Increased CPU utilization to %d%% for %d minutes", utilization, minutes),
		})
	}

	return
}

func (ah *AppHandler) StopCPU(context *gin.Context) {
	ah.logger.Info("Stopping CPU utilization")
	cpuWaster, ok := ah.metrics["cpu"].(*CPUWaster)
	if !ok {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "CPU waster not configured properly.",
		})
		return
	}
	err := cpuWaster.StopCPU()
	if err != nil {
		context.JSON(http.StatusConflict, gin.H{
			"message": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Stopped CPU utilization",
	})
}

func (ah *AppHandler) sendMetrics(msg string, context *gin.Context) {

	param, err := strconv.Atoi(context.Params.ByName("metricValue"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("invalid argument metricValue %s", err.Error()),
		})
		return
	}
	metricsValue := float64(param)
	ah.logger.Info("received request", lager.Data{"metricsValue": metricsValue})
	resp, err := SubmitScaleUpEventToAutoscaler(ah.logger, ah.appConfig, metricsValue, ah.metricsServerBaseUrl)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}
	if resp.StatusCode != http.StatusOK {
		defer func() { _ = resp.Body.Close() }()
		context.JSON(resp.StatusCode, gin.H{
			"message": fmt.Sprintf("Autoscaler responded with %d", resp.StatusCode),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s %f", msg, metricsValue),
	})
}

func postScaleUpEventToAutoscaler(logger lager.Logger, appConfig *cfenv.App, metricsValue float64, autoscalerApiServerUrl string) (*http.Response, error) {
	if appConfig == nil {
		return nil, fmt.Errorf("appConfig cannot be empty")
	}

	appId := appConfig.AppID
	_, err := appConfig.Services.WithName("ak-test-autoscaler")
	if err != nil {
		return nil, err
	}
	cfInstanceKey := os.Getenv("CF_INSTANCE_KEY")
	cfInstanceCert := os.Getenv("CF_INSTANCE_CERT")

	cert, err := tls.LoadX509KeyPair(cfInstanceCert, cfInstanceKey)
	if err != nil {
		log.Printf("Error creating x509 keypair from client cert file %s and client key file %s", cfInstanceCert, cfInstanceKey)
		logger.Error("unable to load x509 keypair", err)
	}
	log.Printf("CAFile: %s", cfInstanceCert)

	caCertBytes, err := os.ReadFile(cfInstanceCert)
	if err != nil {
		log.Printf("Error opening cert file %s, Error: %s", caCertBytes, err)
		logger.Error("unable to read CFInstanceCert keypair", err, lager.Data{cfInstanceCert: cfInstanceCert})
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBytes)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            caCertPool,
	}

	t := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: t}

	metricsValueBody := createCustomMetricsPayload(appId, metricsValue)

	resp, _ := sendRequestToAutoscaler(logger, appId, client, autoscalerApiServerUrl, metricsValueBody, cfInstanceCert)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func sendRequestToAutoscaler(logger lager.Logger, appId string, client *http.Client, autoscalerApiServerUrl string, metricsValueBody []byte, cfInstanceCert string) (*http.Response, error) {
	logger.Info("sending POST to autoscaler")
	customMetricsURL := autoscalerApiServerUrl + "/v1/apps/" + appId + "/metrics"
	logger.Info("sending POST to autoscaler", lager.Data{"autoscalerURL": customMetricsURL})

	log.Printf("custom metrics body: %s ", string(metricsValueBody))
	request, _ := http.NewRequest("POST", customMetricsURL, bytes.NewReader(metricsValueBody))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-Forwarded-Client-Cert", mustReadXFCCcert(cfInstanceCert))
	resp, err := client.Do(request)
	if err != nil {
		logger.Error("failed sending POST to autoscaler", err)
		return nil, fmt.Errorf("unable to send %s %w", request.URL, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return resp, nil
}

func mustReadXFCCcert(fileName string) string {
	file, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading file %s", fileName)
	}
	block, _ := pem.Decode(file)
	return base64.StdEncoding.EncodeToString(block.Bytes)
}
