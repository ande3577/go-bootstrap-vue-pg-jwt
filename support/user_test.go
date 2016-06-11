package support_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("user", func() {
	var u *model.MockUser
	var s *model.MockSession
	var password string
	var passwordConfirmation string

	BeforeEach(func() {
		u = &model.MockUser{User: model.User{Id: 12345, Login: "login", Email: "user@mail.com"}}

		password = "password"
		passwordConfirmation = "password"

		s = &model.MockSession{}
	})

	var _ = Describe("Creating a new user", func() {

		Context("with valid user", func() {
			It("should create a new user", func() {
				By("creating the user", func() {
					err := support.CreateUser(u, s, password, passwordConfirmation)
					Expect(err).To(BeNil())
					Expect(u.Created).To(BeTrue())

					Expect(u.HashedPassword).ToNot(HaveLen(0))
				})

				By("creating a session for the user", func() {
					Expect(s.CreateCalled).To(BeTrue())
					Expect(s.UserId).To(Equal(u.Id))
				})
			})
		})

		Context("with invalid user", func() {
			AfterEach(func() {
				err := support.CreateUser(u, s, password, passwordConfirmation)
				Expect(err).To(HaveOccurred())
				Expect(u.Created).To(BeFalse())
			})

			It("should return an error if login missing", func() {
				u.Login = ""
			})

			It("should return an error if email missing", func() {
				u.Email = ""
			})

			It("should return an error if password missing", func() {
				password = ""
				passwordConfirmation = ""
			})

			It("should return an error if confirmation does not match", func() {
				passwordConfirmation = "different_password"
			})
		})
	})

	var _ = Describe("Updating a user", func() {
		It("should call update", func() {
			passwordChanged, _, err := support.UpdateUser(u, s, "", "", true, true)
			Expect(err).ToNot(HaveOccurred())
			Expect(u.Updated).To(BeTrue())
			Expect(passwordChanged).To(BeFalse())
			Expect(u.SessionsDestroyed).To(BeFalse())
		})

		It("should update password", func() {
			passwordChanged, tokenData, err := support.UpdateUser(u, s, "new_password", "new_password", true, true)
			Expect(err).ToNot(HaveOccurred())
			Expect(u.Updated).To(BeTrue())
			Expect(passwordChanged).To(BeTrue())
			Expect(u.SessionsDestroyed).To(BeTrue())
			Expect(tokenData.TokenString).ToNot(HaveLen(0))
		})

		It("should require passwords to match", func() {
			passwordChanged, _, err := support.UpdateUser(u, s, "new_password", "mismatched_password", true, true)
			Expect(err).To(HaveOccurred())
			Expect(u.Updated).To(BeFalse())
			Expect(passwordChanged).To(BeFalse())
			Expect(u.SessionsDestroyed).To(BeFalse())
		})
	})
})
