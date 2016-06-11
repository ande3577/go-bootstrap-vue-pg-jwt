package integration_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Register", func() {
	Context("as unauthenticated user", func() {

		BeforeEach(func() {
			logout()
			Expect(page.FindByLink("Register").Click()).To(Succeed())
			Eventually(page).Should(HaveURL("http://localhost:5000/register"))

			Expect(page.Find("#login").Fill("user")).Should(Succeed())
			Expect(page.Find("#email").Fill("user@mail.com")).Should(Succeed())

			Expect(page.Find("#password-main").Fill("password")).Should(Succeed())
			Expect(page.Find("#password-confirmation").Fill("password")).Should(Succeed())
		})

		AfterEach(func() {
			var u *model.User = &model.User{}
			if err := u.FindByLogin("user"); err == nil {
				if err = u.Destroy(); err != nil {
					panic(err)
				}
			}
		})

		It("should allow restering a new user", func() {
			By("registering and redirecting to the main page", func() {
				Expect(page.FindByButton("Register").Click()).To(Succeed())
				Eventually(page).Should(HaveURL("http://localhost:5000/"))
			})

			By("hiding login link after registered", func() {
				Expect(page.FindByButton("Login")).ToNot(BeFound())
			})

			By("showing user name after user is registered", func() {
				Expect(getCurrentlyLoggedInUser()).To(Equal("user"))
			})
		})

		It("should show error for invalid registrations", func() {
			Expect(page.Find("#password-confirmation").Fill("different_password")).Should(Succeed())
			Expect(page.FindByButton("Register").Click()).To(Succeed())
			Eventually(page).Should(HaveURL("http://localhost:5000/register")) // redirect back to register
			ExpectErrorMessage()
			Expect(page.Find("#login")).To(HaveValue("user"))          // preserve login
			Expect(page.Find("#email")).To(HaveValue("user@mail.com")) // preserve email
		})

	})
})
