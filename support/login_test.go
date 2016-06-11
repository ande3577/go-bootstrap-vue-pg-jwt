package support_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/model"
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/support"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Logging in for a user", func() {
	Context("user exists", func() {
		var u *model.MockUser
		var s *model.MockSession

		BeforeEach(func() {
			u = &model.MockUser{User: model.User{Login: "username", HashedPassword: auth.GenerateHashFromPassword("password"), Email: "user@mail.com"}}
			s = &model.MockSession{}
		})

		It("should login", func() {
			tokenData, err := support.Login("username", "password", true, false, u, s)
			Expect(err).To(BeNil())
			Expect(tokenData.UserId).To(Equal("username"))
			Expect(s.CreateCalled).To(BeTrue())
			Expect(s.Session.Session).To(Equal(tokenData.SessionIdentifier))
		})

		It("should not login for incorrect password", func() {
			_, err := support.Login("username", "wrong_password", true, false, u, s)
			Expect(err).To(HaveOccurred())
			Expect(s.CreateCalled).To(BeFalse())
		})

		It("should login by email", func() {
			tokenData, err := support.Login(u.Email, "password", true, false, u, s)
			Expect(err).To(BeNil())
			Expect(tokenData.UserId).To(Equal("username"))
		})

	})
})

var _ = Describe("Logging out for a user", func() {
	Context("with a logged in user", func() {
		var u *model.MockUser
		var s *model.MockSession

		BeforeEach(func() {
			u = &model.MockUser{User: model.User{Id: 1234, Login: "username", HashedPassword: auth.GenerateHashFromPassword("password"), Email: "user@mail.com"}}
			s = &model.MockSession{Session: model.Session{UserId: u.Id, Session: "asdfasdf"}}
		})

		It("should destroy active session", func() {
			support.Logout(u.Login, s.Session.Session, u, s)
			Expect(s.DestroyCalled).To(BeTrue())
		})
	})
})
