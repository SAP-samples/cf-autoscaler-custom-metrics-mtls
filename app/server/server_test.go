package server_test

import (
	"github.com/asalan316/golang-autoscaler-custom-metrics/app/server"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
)

var _ = Describe("Server Tests", func() {
	var appConfig *cfenv.App

	Context("New Server is called", func() {
		It("return a new server", func() {
			validEnv := []string{
				`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
				`VCAP_SERVICES={"autoscaler":[{"binding_guid":"eeb9c732-69f4-4d2e-b7f570f","binding_name":null,"tags": [ "autoscaler", "app-autoscaler","cf-autoscaler"],"instance_guid":"5f6e545e-c5cf-42a6-86be-64778b9d6a86","instance_name":"ak-test-autoscaler","label":"autoscaler","name":"ak-test-autoscaler","plan":"standard","provider":null,"syslog_drain_url":null,"volume_mounts":[]}]}`,
			}
			testEnv := cfenv.Env(validEnv)
			appConfig, err := cfenv.New(testEnv)
			Expect(err).Should(Not(HaveOccurred()))
			server := server.NewServer(appConfig)

			Expect(server).NotTo(BeNil())
		})

	})

	Context("/ endpoint is called", func() {
		BeforeEach(func() {
			validEnv := []string{
				`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
				`VCAP_SERVICES={"autoscaler": [{"Name": "ak-test-autoscaler","Label": "autoscaler", "Tags": [ "autoscaler","app-autoscaler","cf-autoscaler"],"Plan": "standard","Credentials": null,"VolumeMounts": null }]}`,
			}
			testEnv := cfenv.Env(validEnv)
			var err error
			appConfig, err = cfenv.New(testEnv)
			Expect(err).Should(Not(HaveOccurred()))

			os.Setenv("CF_INSTANCE_CERT", "api_public.crt")
		})

		It("returns configurations with HTTP 200", func() {
			mockResponse := `{
        "CF_INSTANCE_CERT": "api_public.crt",
        "CF_INSTANCE_KEY": "",
        "appConfigurations": {
          "instance_id": "451f045fd16427bb99c895a2649b7b2a",
          "application_id": "abcabc123123defdef456456",
          "instance_index": 0,
          "name": "styx-james",
          "host": "0.0.0.0",
          "port": 61857,
          "version": "c1063c1c-40b9-434e-a797-db240b587d32",
          "application_uris": [
            "styx-james.a1-app.cf-app.com"
          ],
          "space_id": "3e0c28c5-6d9c-436b-b9ee-1f4326e54d05",
          "space_name": "jdk",
          "Home": "",
          "MemoryLimit": "",
          "WorkingDir": "",
          "TempDir": "",
          "User": "",
          "Services": {
            "autoscaler": [
              {
                "Name": "ak-test-autoscaler",
                "Label": "autoscaler",
                "Tags": [
                  "autoscaler",
                  "app-autoscaler",
                  "cf-autoscaler"
                ],
                "Plan": "standard",
                "Credentials": null,
                "VolumeMounts": null
              }
            ]
          },
          "cf_api": "https://api.system_domain.com",
          "limits": {
            "disk": 1024,
            "fds": 16384,
            "mem": 512
          }
        }}`
			router := setUpRouter()
			ah := server.NewAppHandler(appConfig, "")
			router.GET("/", ah.GetHome)

			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseData, _ := io.ReadAll(w.Body)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(string(responseData)).To(MatchJSON(mockResponse))
		})

	})

	Context("/busy endpoint", func() {

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

		It("/busy Endpoint is called", func() {
			mockResponse := `{
        				"message": "I am busy with value 300.000000"
      				}`
			router := setUpRouter()
			ah := server.NewAppHandler(appConfig, autoscalerApiServer.URL())
			router.GET("/busy/:metricValue", ah.Busy)

			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.AddParam("metricValue", "199")
			/*// mocking a function - TODo - may be use api call - //TODO consider httest server https://speedscale.com/testing-golang-with-httptest/

			 */

			// TODO Monkey patching conflicts with responses
			/*server.SubmitScaleUpEventToAutoscaler = func(appConfig *cfenv.App, metricsValue float64, apiServerUrl string) (*http.Response, error) {
				return &http.Response{
					Status:     "",
					StatusCode: 200,
					Body:       nil,
				}, nil
			}*/
			// that's why spin a server
			autoscalerApiServer.RouteToHandler(
				http.MethodPost, "/v1/apps/abcabc123123defdef456456/metrics",
				ghttp.RespondWithJSONEncoded(http.StatusOK,
					nil))

			req, _ := http.NewRequest(http.MethodGet, "/busy/300", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseData, _ := ioutil.ReadAll(w.Body)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(string(responseData)).To(MatchJSON(mockResponse))
		})

	})

	Context("/not-busy endpoint", func() {

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

		It("/not-busy Endpoint is called", func() {
			router := setUpRouter()
			ah := server.NewAppHandler(appConfig, autoscalerApiServer.URL())
			router.GET("/not-busy/:metricValue", ah.NotBusy)

			/*server.SubmitScaleUpEventToAutoscaler = func(
				appConfig *cfenv.App, metricsValue float64, apiServerUrl string) (*http.Response, error) {
				return &http.Response{
					Status:     "",
					StatusCode: 200,
					Body:       nil,
				}, nil
			}*/

			autoscalerApiServer.RouteToHandler(
				http.MethodPost, "/v1/apps/abcabc123123defdef456456/metrics",
				ghttp.RespondWithJSONEncoded(http.StatusOK,
					nil))

			req, _ := http.NewRequest(http.MethodGet, "/not-busy/199", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			responseData, _ := io.ReadAll(w.Body)

			Expect(w.Code).To(Equal(http.StatusOK))
			Expect(string(responseData)).To(MatchJSON(`{"message":"I am not busy with value 199.000000"}`))
		})

	})
})

func setUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
