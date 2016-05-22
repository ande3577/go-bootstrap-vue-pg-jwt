package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("Index", func() {
	Context("as unauthenticated user", func() {
		BeforeEach(func() {
			logout()
		})

		It("should not show the main page", func() {
			Eventually(page).Should(HaveURL("http://localhost:5000/"))
		})
	})
})
