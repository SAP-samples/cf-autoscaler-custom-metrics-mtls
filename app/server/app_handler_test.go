package server_test

import (
	"github.com/asalan316/golang-autoscaler-custom-metrics/app/server"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"net/http/httptest"
	"os"
)

var _ = Describe("App Handler Tests", func() {
	var appConfig *cfenv.App

	Context("busy is called", func() {
		BeforeEach(func() {
			validEnv := []string{
				`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
				`VCAP_SERVICES={"autoscaler":[{"binding_guid":"eeb9c732-69f4-4d2e-b7f570f","binding_name":null,"tags": [ "autoscaler", "app-autoscaler","cf-autoscaler"],"instance_guid":"5f6e545e-c5cf-42a6-86be-64778b9d6a86","instance_name":"ak-test-autoscaler","label":"autoscaler","name":"ak-test-autoscaler","plan":"standard","provider":null,"syslog_drain_url":null,"volume_mounts":[]}]}`,
			}
			testEnv := cfenv.Env(validEnv)
			var err error
			appConfig, err = cfenv.New(testEnv)
			Expect(err).Should(Not(HaveOccurred()))

			os.Setenv("CF_INSTANCE_CERT", "api_public.crt")
			autoscalerApiServer = ghttp.NewServer()
		})

		It("with no request param and returns successful, HTTP 400", func() {
			ah := server.NewAppHandler(appConfig, "")
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ah.Busy(ctx)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
		})

		It("with invalid response from autoscaler service and returns InternalServerError, HTTP 500", func() {

			ah := server.NewAppHandler(nil, "")
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.AddParam("metricValue", "300")
			ah.Busy(ctx)
			Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			Expect(responseRecorder.Body).To(MatchJSON(`{"message":"appConfig cannot be empty"}`))
		})

		It("with response 400 from autoscaler service and returns badRequest, HTTP 400", func() {
			autoscalerApiServer.RouteToHandler(
				http.MethodPost, "/v1/apps/abcabc123123defdef456456/metrics",
				ghttp.RespondWithJSONEncoded(http.StatusBadRequest,
					nil))

			ah := server.NewAppHandler(appConfig, autoscalerApiServer.URL())
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.AddParam("metricValue", "300")
			ah.Busy(ctx)
			Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			Expect(responseRecorder.Body).To(MatchJSON(`{"message":"Autoscaler responded with 400"}`))
		})

		It("with correct request param and returns successful, HTTP 200", func() {
			autoscalerApiServer.RouteToHandler(
				http.MethodPost, "/v1/apps/abcabc123123defdef456456/metrics",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/apps/abcabc123123defdef456456/metrics"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, nil)))

			ah := server.NewAppHandler(appConfig, autoscalerApiServer.URL())
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.AddParam("metricValue", "300")

			ah.Busy(ctx)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			Expect(responseRecorder.Body).To(MatchJSON(`{"message":"I am busy with value 300.000000"}`))
		})

	})

	Context("not busy is called", func() {
		BeforeEach(func() {
			validEnv := []string{
				`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"an-app-id","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
				`VCAP_SERVICES={"autoscaler":[{"binding_guid":"eeb9c732-69f4-4d2e-b7f570f","binding_name":null,"tags": [ "autoscaler", "app-autoscaler","cf-autoscaler"],"instance_guid":"5f6e545e-c5cf-42a6-86be-64778b9d6a86","instance_name":"ak-test-autoscaler","label":"autoscaler","name":"ak-test-autoscaler","plan":"standard","provider":null,"syslog_drain_url":null,"volume_mounts":[]}]}`,
			}
			testEnv := cfenv.Env(validEnv)
			var err error
			appConfig, err = cfenv.New(testEnv)
			Expect(err).Should(Not(HaveOccurred()))

			os.Setenv("CF_INSTANCE_CERT", "api_public.crt")
			autoscalerApiServer = ghttp.NewServer()

		})

		It("with correct request param and returns successful, HTTP 200", func() {

			autoscalerApiServer.RouteToHandler(
				http.MethodPost, "/v1/apps/an-app-id/metrics",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/apps/an-app-id/metrics"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, nil)))

			ah := server.NewAppHandler(appConfig, autoscalerApiServer.URL())
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.AddParam("metricValue", "199")

			ah.NotBusy(ctx)

			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			Expect(responseRecorder.Body).To(MatchJSON(`{"message":"I am not busy with value 199.000000"}`))
		})

	})
})
