package app_config_test

import (
	reporter "github.com/joefitzgerald/rainbow-reporter"
	"github.com/sclevine/spec"

	"testing"
)

var suite spec.Suite

func TestSuite(t *testing.T) {
	suite.Run(t)
}

func init() {
	suite = spec.New("config", spec.Report(reporter.Rainbow{}))
	suite("cfenv", testGetAppEnv)
}
