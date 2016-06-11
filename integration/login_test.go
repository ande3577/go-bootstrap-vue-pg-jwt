package integration_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Login", func() {
	Context("as unauthenticated user", func() {
		var u *model.User = &model.User{Login: "user", Email: "user@mail.com"}

		BeforeEach(func() {
			logout()
			s := &model.MockSession{} // don't actually want to create a session here
			err := support.CreateUser(u, s, "password", "password")
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			u.Destroy()
		})

		It("should allow logging in", func() {
			By("logging in and redirecting to main page", func() {
				loginAs("user", "password")
			})

			By("hiding login link after user is logged in", func() {
				Expect(page.FindByButton("Login")).ToNot(BeFound())
			})

			By("showing user name after user is logged in", func() {
				Expect(getCurrentlyLoggedInUser()).To(Equal("user"))
			})

			By("allowing user to logout", func() {
				Expect(page.FindByButton("Logout")).To(BeFound())
				Expect(page.FindByButton("Logout").Click()).To(Succeed())
				Eventually(page).Should(HaveURL("http://localhost:5000/"))
				Expect(page.FindByButton("Login")).To(BeFound())
			})

		})

		It("should allow sign in by email", func() {
			By("redirecting to the index page", func() {
				loginAs(u.Email, "password")
				Eventually(page).Should(HaveURL("http://localhost:5000/"))
			})

			By("displaying users name", func() {
				Expect(getCurrentlyLoggedInUser()).To(Equal("user"))
			})
		})

		It("should show error for invalid login", func() {
			Eventually(page.Find("#password")).Should(BeFound())
			Expect(page.Find("#password").Fill("invalid")).Should(Succeed())
			Expect(page.FindByButton("Login").Click()).Should(Succeed())
			Eventually(page).Should(HaveURL("http://localhost:5000/login"))
			Eventually(page.First("#error-message")).Should(BeFound())
			loginAs("user", "password")
		})

	})
})
