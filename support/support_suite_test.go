package support_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"flag"
	"github.com/onsi/ginkgo/reporters"
	"testing"
)

var skipIntegration = flag.Bool("skip-integration", false, "skip all integration tests")
var leaveBrowserOpen = flag.Bool("leave-browser-open", false, "leave browser open after integration tests")

func TestSupport(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Support Suite", []Reporter{junitReporter})
}
