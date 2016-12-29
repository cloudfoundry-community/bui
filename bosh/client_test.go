package bosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-community/bui/bosh"
)

var _ = Describe("BOSH Client", func() {
	Describe("Test Default Config", func() {
		config := DefaultConfig()

		It("returns default config", func() {
			Expect(config.BOSHAddress).Should(Equal("https://192.168.50.4:25555"))
			Expect(config.SkipSslValidation).Should(Equal(true))
		})

	})

	Describe("Test Creating client", func() {
		var client *Client

		BeforeEach(func() {
			setup(MockRoute{"GET", "/stemcells", `{}`, ""}, "basic")
			config := &Config{
				BOSHAddress: server.URL,
			}

			client, _ = NewClient(config)
		})

		AfterEach(func() {
			teardown()
		})

		It("can get bosh info", func() {
			info, err := client.GetInfo()
			Expect(info.Name).Should(Equal("bosh-lite"))
			Expect(info.UUID).Should(Equal("2daf673a-9755-4b4f-aa6d-3632fbed8019"))
			Expect(info.Version).Should(Equal("1.3126.0 (00000000)"))
			Expect(info.User).Should(Equal("admin"))
			Expect(info.CPI).Should(Equal("warden_cpi"))

			Expect(err).Should(BeNil())
		})
	})
})
