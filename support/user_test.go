package support_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Creating new user", func() {
	var u *model.MockUser
	var password string
	var passwordConfirmation string

	BeforeEach(func() {
		u = &model.MockUser{User: model.User{Login: "login", Email: "user@mail.com"}}
		u.Login = "login"
		u.Email = "user@mail.com"
		password = "password"
		passwordConfirmation = "password"
	})

	Context("with valid user", func() {
		It("should create a new user", func() {
			err := support.CreateUser(u, password, passwordConfirmation)
			Expect(err).To(BeNil())
			Expect(u.Created).To(BeTrue())

			Expect(u.HashedPassword).ToNot(HaveLen(0))
		})
	})

	Context("with invalid user", func() {
		AfterEach(func() {
			err := support.CreateUser(u, password, passwordConfirmation)
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

var _ = Describe("Logging in for a user", func() {
	Context("user exists", func() {
		var u *model.MockUser

		BeforeEach(func() {
			u = &model.MockUser{User: model.User{Login: "username", HashedPassword: auth.GenerateHashFromPassword("password"), Email: "user@mail.com"}}
		})

		It("should login", func() {
			login, err := support.Login("username", "password", u)
			Expect(err).To(BeNil())
			Expect(login).To(Equal("username"))
		})

		It("should login by email", func() {
			login, err := support.Login(u.Email, "password", u)
			Expect(err).To(BeNil())
			Expect(login).To(Equal("username"))
		})

	})
})
