package app_config_test

import (
	"fmt"
	appconfig "github.com/asalan316/golang-autoscaler-custom-metrics/app/app_config"
	"github.com/cloudfoundry-community/go-cfenv"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"os"
	"testing"
)

func testGetAppEnv(t *testing.T, when spec.G, it spec.S) {
	it.Before(func() {
		RegisterTestingT(t)
	})
	validEnv := []string{
		`VCAP_APPLICATION={"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`,
		`VCAP_SERVICES={"autoscaler":[{"binding_guid":"eeb9c732-69f4-4d2e-b7f570f","binding_name":null,"tags": [ "autoscaler", "app-autoscaler","cf-autoscaler"],"instance_guid":"5f6e545e-c5cf-42a6-86be-64778b9d6a86","instance_name":"ak-test-autoscaler","label":"autoscaler","name":"ak-test-autoscaler","plan":"standard","provider":null,"syslog_drain_url":null,"volume_mounts":[]}]}`,
	}
	// Mocking Technique#1 - Higher order function - use to mock the package level functions
	currentAppEnv := func() (*cfenv.App, error) {
		testEnv := cfenv.Env(validEnv)
		return cfenv.New(testEnv)
	}

	when("getAppEnv is invoked", func() {
		it.Before(func() {
			os.Setenv("VCAP_APPLICATION", "{}")
		})
		it("returns vcap services configuration", func() {
			appEnv, err := appconfig.GetAppEnv(currentAppEnv)
			Expect(err).Should(Not(HaveOccurred()))
			fmt.Println(appEnv)
			Expect(appEnv.Services["autoscaler"][0].Name).To(Equal("ak-test-autoscaler"))
			Expect(appEnv.Services["autoscaler"][0].Label).To(Equal("autoscaler"))
			Expect(appEnv.Services["autoscaler"][0].Tags).To(Equal([]string{"autoscaler", "app-autoscaler", "cf-autoscaler"}))

		})
	})

	when("getAppEnv is invoked and not running on cf", func() {
		it.Before(func() {
			os.Setenv("VCAP_APPLICATION", "")
		})
		it("IsRunningOnCF() returns false", func() {
			/* Mocking Technique#2
			Monkey Patching - use to mock the package level functions - cfenv.IsRunningOnCF
			*/
			appconfig.IsRunningOnCF = func() bool {
				return false
			}
			_, err := appconfig.GetAppEnv(currentAppEnv)
			Expect(err).To(HaveOccurred())
		})
	})

}
