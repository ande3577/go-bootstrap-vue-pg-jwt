package auth_test

import (
	"github.com/ande3577/go-bootstrap-vue-pg-jwt/auth"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	BeforeSuite(func() {
		cookieStoreAuthenticationKey, _ := auth.GenerateRandomString(32)
		cookieStoreEncryptionKey, _ := auth.GenerateRandomString(32)
		jwtSigningKeyString, _ := auth.GenerateRandomString(32)
		cookieIssuer := "auth-test"

		auth.Initialize(cookieStoreAuthenticationKey, cookieStoreEncryptionKey, jwtSigningKeyString, cookieIssuer)
	})

	Describe("Login and parse token", func() {
		It("should create an auth token", func() {
			tokenData, err := auth.Login("user", false, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenData.UserId).To(Equal("user"))
			Expect(tokenData.XsrfToken).ToNot(HaveLen(0))
			_, err = auth.ParseToken(tokenData.TokenString, false, false)
			Expect(err).ToNot(HaveOccurred())
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
	})

})
