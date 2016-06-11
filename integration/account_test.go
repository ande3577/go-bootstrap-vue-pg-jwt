package integration_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("account settings", func() {
	Context("with a logged in user", func() {
		var u *model.User

		BeforeEach(func() {
			logout()
			u = createUser()
			loginAsUser()

			Expect(page.FindByLink("user")).To(BeFound())
			Expect(page.FindByLink("user").Click()).To(Succeed())
			Eventually(page.Find("#email")).Should(BeFound())
		})

		AfterEach(func() {
			u.Destroy()
		})

		It("should allow visiting account page", func() {
			Expect(page.Find("#email")).To(HaveValue("user@mail.com"), "load email")
			Expect(page.Find("#password-main")).To(HaveValue(""))
		})

		It("should allow changing email", func() {
			By("editing email", func() {
				Expect(page.Find("#email").Fill("new_user@new_mail.com")).To(Succeed())
			})

			By("submitting form", func() {
				Expect(page.FindByButton("Submit").Click()).To(Succeed())
				Eventually(page).Should(HaveURL("http://localhost:5000/"), "redirect to main")
			})

			By("changing email", func() {
				Expect(page.FindByLink("user").Click()).To(Succeed())
				Eventually(page.Find("#email")).Should(BeFound())
				Expect(page.Find("#email")).To(HaveValue("new_user@new_mail.com"), "load new email")
			})
		})

		It("should allow changing password", func() {
			// create a session to ensure that is destroyed when changing the password
			s := &model.Session{UserId: u.Id, Session: "asdfasdfasdf"}
			err := s.Create()
			Expect(err).ToNot(HaveOccurred())

			By("submitting form", func() {
				Expect(page.Find("#password-main").Fill("new_password")).To(Succeed())
				Expect(page.Find("#password-confirmation").Fill("new_password")).To(Succeed())
				Expect(page.FindByButton("Submit").Click()).To(Succeed())
			})

			By("remaining logged in", func() {
				Eventually(page).Should(HaveURL("http://localhost:5000/"), "redirect to main")
				Expect(getCurrentlyLoggedInUser()).To(Equal("user"))
			})

			By("still having an active session", func() {
				Expect(page.FindByLink("user")).To(BeFound())
				Expect(page.FindByLink("user").Click()).To(Succeed())
				Eventually(page.Find("#email")).Should(BeFound())
			})

			By("changing password", func() {
				logout()
				loginAs("user", "new_password")
				Eventually(getCurrentlyLoggedInUser()).Should(Equal("user"))
			})

			By("invalidating old sessions", func() {
				err := s.FindBySession(s.Session)
				Expect(err).To(HaveOccurred(), "should not be able to destroy this session, since it should already have been destroyed")
			})
		})

		It("should require passwords to match", func() {
			By("remaining on account page", func() {
				Expect(page.Find("#email").Fill("new_user@new_mail.com")).To(Succeed()) // change the email as well
				Expect(page.Find("#password-main").Fill("new_password")).To(Succeed())
				Expect(page.Find("#password-confirmation").Fill("mismatched_password")).To(Succeed())
				Expect(page.FindByButton("Submit").Click()).To(Succeed())

				Eventually(page.Find("#email")).Should(BeFound())
			})

			By("displaying error message", func() {
				ExpectErrorMessage()
			})

			By("leaving email entered", func() {
				Expect(page.Find("#email")).To(HaveValue("new_user@new_mail.com"), "load new email")
			})
		})

		It("get should deny access if not signed in", func() {
			logout()
			page.Navigate("http://localhost:5000/account")
			ExpectErrorMessage()
		})

		Context("without valid session", func() {
			BeforeEach(func() {
				err := u.DestroySessions()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should validate session on get", func() {
				Expect(page.FindByLink("user")).To(BeFound())
				Expect(page.FindByLink("user").Click()).To(Succeed())
				ExpectErrorMessage()
			})

			It("should validate session on post", func() {
				Expect(page.FindByButton("Submit").Click()).To(Succeed())
				ExpectErrorMessage()
			})
		})
	})
})
