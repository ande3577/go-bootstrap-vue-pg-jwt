package auth_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"

	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockSessionProvider struct {
	session    string
	findCalled bool
}

func (s *MockSessionProvider) FindBySession(session string) error {
	s.findCalled = true
	if s.session == session {
		return nil
	} else {
		return errors.New("session not found")
	}
}

var _ = Describe("Auth", func() {
	var sessionProvider *MockSessionProvider

	BeforeSuite(func() {
		cookieStoreAuthenticationKey, _ := auth.GenerateRandomString(32)
		cookieStoreEncryptionKey, _ := auth.GenerateRandomString(32)
		jwtSigningKeyString, _ := auth.GenerateRandomString(32)
		cookieIssuer := "auth-test"
		sessionProvider = &MockSessionProvider{}

		auth.Initialize(sessionProvider, cookieStoreAuthenticationKey, cookieStoreEncryptionKey, jwtSigningKeyString, cookieIssuer)
	})

	Describe("Login and parse token", func() {
		It("should create an auth token", func() {
			tokenData, err := auth.Login("user", false, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenData.UserId).To(Equal("user"))
			Expect(tokenData.XsrfToken).ToNot(HaveLen(0))
			Expect(tokenData.SessionIdentifier).ToNot(HaveLen(0))
			Expect(tokenData.TokenString).ToNot(HaveLen(0), "create token")
			Expect(tokenData.RefreshTokenString).To(HaveLen(0), "do not create refresh token")

			tokenDataOut, err := auth.ParseToken(tokenData.TokenString, false, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(tokenDataOut.UserId).To(Equal(tokenData.UserId))
			Expect(tokenDataOut.XsrfToken).To(Equal(tokenData.XsrfToken))
			Expect(tokenDataOut.SessionIdentifier).To(Equal(tokenData.SessionIdentifier))
		})

		It("should create an http auth token", func() {
			tokenData, err := auth.Login("user", true, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(tokenData.TokenString).ToNot(HaveLen(0), "create token")
			Expect(tokenData.RefreshTokenString).ToNot(HaveLen(0), "create refresh token")
			Expect(tokenData.TokenString).ToNot(Equal(tokenData.RefreshTokenString), "ensure the two tokens are not equal")
		})

		It("should reject invalid tokens", func() {
			_, err := auth.ParseToken("asdfasdfsdf", false, false)
			Expect(err).To(HaveOccurred())
		})

		It("should reject tokens created in development mode", func() {
			tokenData, err := auth.Login("user", false, true)
			Expect(err).ToNot(HaveOccurred())
			_, err = auth.ParseToken(tokenData.TokenString, false, false)
			Expect(err).To(HaveOccurred())
		})

		It("should reject tokens created in http request for json request", func() {
			tokenData, err := auth.Login("user", true, false)
			Expect(err).ToNot(HaveOccurred())
			_, err = auth.ParseToken(tokenData.TokenString, false, false)
			Expect(err).To(HaveOccurred())
		})

		It("should not assign a session id if user id is blank", func() {
			tokenData, err := auth.Login("", true, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenData.SessionIdentifier).To(HaveLen(0))
		})
	})

	Describe("handle refresh token", func() {
		It("should handle not expired token", func() {
			tokenData, err := auth.Login("user", true, false)
			Expect(err).ToNot(HaveOccurred(), "login")
			newToken, _, err := auth.ParseTokenWithRefreshToken(tokenData.TokenString, tokenData.RefreshTokenString, true, false)
			Expect(err).ToNot(HaveOccurred(), "parse token does not create error")
			Expect(sessionProvider.findCalled).To(BeFalse())
			Expect(newToken).To(BeFalse())
		})

		It("should check the refresh token to obtain a new token", func() {
			tokenData, err := auth.Login("user", true, false)
			tokenData.TokenString = "" // clear out the token string
			Expect(err).ToNot(HaveOccurred(), "login")
			sessionProvider.session = tokenData.SessionIdentifier
			newToken, tokenData, err := auth.ParseTokenWithRefreshToken(tokenData.TokenString, tokenData.RefreshTokenString, true, false)
			Expect(err).ToNot(HaveOccurred(), "parse token does not create error")
			Expect(sessionProvider.findCalled).To(BeTrue(), "checked the session store")
			Expect(tokenData.TokenString).ToNot(HaveLen(0), "should assign a new token")
			Expect(tokenData.TokenString).ToNot(Equal(tokenData.RefreshTokenString), "ensure the two tokens are not equal")
			Expect(newToken).To(BeTrue())
			Expect(tokenData.UserId).To(Equal("user"))
		})

		It("should not issue a new token if the session cannot be validated", func() {
			tokenData, err := auth.Login("user", true, false)
			Expect(err).ToNot(HaveOccurred(), "login")
			sessionProvider.session = ""
			_, _, err = auth.ParseTokenWithRefreshToken("", tokenData.RefreshTokenString, true, false)
			Expect(err).To(HaveOccurred())
		})
	})
})
