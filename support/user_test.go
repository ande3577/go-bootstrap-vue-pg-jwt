package support_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Creating new user", func() {
	var u *model.MockUser
	var s *model.MockSession
	var password string
	var passwordConfirmation string

	BeforeEach(func() {
		u = &model.MockUser{User: model.User{Login: "login", Email: "user@mail.com"}}
		u.Login = "login"
		u.Email = "user@mail.com"
		u.Id = 12345
		password = "password"
		passwordConfirmation = "password"

		s = &model.MockSession{}
	})

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
