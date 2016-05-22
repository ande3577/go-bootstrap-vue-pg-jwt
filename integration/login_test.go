package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Login", func() {
	Context("as unauthenticated user", func() {
		BeforeEach(func() {
			logout()
		})

		It("should allow logging in", func() {
			By("logging in and redirecting to main page", func() {
				loginAs("user", "user")
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

		It("should show error for invalid login", func() {
			Eventually(page.Find("#password")).Should(BeFound())
			Expect(page.Find("#password").Fill("invalid")).Should(Succeed())
			Expect(page.FindByButton("Login").Click()).Should(Succeed())
			Eventually(page).Should(HaveURL("http://localhost:5000/login"))
			Eventually(page.First("#error-message")).Should(BeFound())
			loginAs("user", "user")
		})

	})
})
