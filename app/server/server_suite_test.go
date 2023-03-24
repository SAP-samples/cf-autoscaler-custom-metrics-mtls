package server_test

import (
	"github.com/onsi/gomega/ghttp"
	"testing"
)

import . "github.com/onsi/gomega"
import . "github.com/onsi/ginkgo/v2"

var autoscalerApiServer *ghttp.Server

func TestApp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "server Suite")
}
