package integration_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/app"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/routes"

	"flag"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
	"testing"
)

var queryOverridePartNumber string
var queryOverrideVersion string

var skipIntegration = flag.Bool("skip-integration", false, "skip all integration tests")
var leaveBrowserOpen = flag.Bool("leave-browser-open", false, "leave browser open after integration tests")

func TestIntegration(t *testing.T) {
	if *skipIntegration {
		t.Skip()
	}

	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Integration Suite", []Reporter{junitReporter})
}

var agoutiDriver *agouti.WebDriver
var page *agouti.Page

func getCurrentlyLoggedInUser() string {
	if page.Find("#current-user") != nil {
		if userId, err := page.Find("#current-user").Text(); err == nil {
			return userId
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func logout() {
	Expect(page.Navigate("http://localhost:5000")).To(Succeed())
	Eventually(page).Should(HaveURL("http://localhost:5000/"))
	page.FindByButton("Logout").Click()
	Eventually(page).Should(HaveURL("http://localhost:5000/"))
	Eventually(page.FindByButton("Login")).Should(BeFound())
}

func loginAs(username string, password string) {
	// if we are already logged in as the target user, do nothing
	if getCurrentlyLoggedInUser() == username {
		return
	}
	logout()
	Eventually(page).Should(HaveURL("http://localhost:5000/"))
	Eventually(page.Find("#user_id")).Should(BeFound())
	Expect(page.Find("#user_id").Fill(username)).Should(Succeed())
	Eventually(page.Find("#password")).Should(BeFound())
	Expect(page.Find("#password").Fill(password)).Should(Succeed())
	Expect(page.FindByButton("Login").Click()).Should(Succeed())
	Eventually(page).Should(HaveURL("http://localhost:5000/"))
}

var appSettings = &app.ApplicationSettings{
	Environment:     "test",
	Port:            "5000",
	RootDirectory:   "./..",
	DevelopmentMode: true,
}

var _ = BeforeSuite(func() {
	// Choose a WebDriver:

	// agoutiDriver = agouti.PhantomJS()
	// agoutiDriver = agouti.Selenium()
	agoutiDriver = agouti.ChromeDriver()

	Expect(agoutiDriver.Start()).To(Succeed())

	app.Initialize(appSettings)
	routes.SetupApplication(appSettings)

	var err error
	page, err = agoutiDriver.NewPage()
	Expect(err).NotTo(HaveOccurred())
	Expect(page.Navigate("http://localhost:5000")).To(Succeed())
})

var _ = AfterSuite(func() {
	if !*leaveBrowserOpen {
		if page != nil {
			Expect(page.Destroy()).To(Succeed())
		}
		Expect(agoutiDriver.Stop()).To(Succeed())
	}
})
