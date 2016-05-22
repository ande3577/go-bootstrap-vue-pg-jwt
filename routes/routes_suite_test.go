package routes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"flag"
	"github.com/onsi/ginkgo/reporters"
	"testing"
)

var skipIntegration = flag.Bool("skip-integration", false, "skip all integration tests")

func TestRoutes(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Routes Suite", []Reporter{junitReporter})
}
