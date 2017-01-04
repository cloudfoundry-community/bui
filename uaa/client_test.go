package uaa_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-community/bui/uaa"
)

var _ = Describe("UAA Client", func() {
	Describe("Test Default Config", func() {
		config := DefaultConfig()

		It("returns default config", func() {
			Expect(config.Address).Should(Equal("https://192.168.50.4:8443"))
			Expect(config.SkipSslValidation).Should(Equal(true))
		})

	})

	Describe("Test Creating client", func() {
		var client *Client

		BeforeEach(func() {
			setup(MockRoute{"GET", "/info", info, ""})
			config := &Config{
				Address: server.URL,
			}

			client, _ = NewClient(config)
		})

		AfterEach(func() {
			teardown()
		})

		It("can get uaa info", func() {
			info, err := client.GetInfo()
			Expect(info.App.Version).Should(Equal("3.4.2"))
			Expect(info.Links.UAA).Should(Equal("https://10.244.66.2:8443"))
			Expect(info.Links.Password).Should(Equal("/forgot_password"))
			Expect(info.Links.Login).Should(Equal("https://10.244.66.2:8443"))
			Expect(info.Links.Register).Should(Equal("/create_account"))
			Expect(info.ZoneName).Should(Equal("uaa"))
			Expect(err).Should(BeNil())
		})
	})

	Describe("Test get password token", func() {
		var client *Client
		BeforeEach(func() {
			setup(MockRoute{"POST", "/oauth/token", tokenResp, ""})

			config := &Config{
				Address:      server.URL,
				ClientID:     "uaa",
				ClientSecret: "uaa-secret",
			}

			client, _ = NewClient(config)
		})

		AfterEach(func() {
			teardown()
		})

		It("can get password token", func() {
			tokenResp, err := client.GetPasswordToken("foo", "bar")
			Expect(err).Should(BeNil())
			Expect(tokenResp.AccessToken).Should(Equal("eyJhbGciOiJSUzI1NiIsImtpZCI6ImxlZ2FjeS10b2tlbi1rZXkiLCJ0eXAiOiJKV1QifQ.eyJqdGkiOiJjNWE0MDA1MjllMDE0NjZkOTdhODE4N2VkMGMyNGU4OSIsInN1YiI6IjBmY2U5ZTc0LTAzMzgtNGZkMC1iMTI1LTVjYjE3NTBjN2FiZCIsInNjb3BlIjpbIm9wZW5pZCIsImJvc2guYWRtaW4iXSwiY2xpZW50X2lkIjoiYm9zaF9jbGkiLCJjaWQiOiJib3NoX2NsaSIsImF6cCI6ImJvc2hfY2xpIiwiZ3JhbnRfdHlwZSI6InBhc3N3b3JkIiwidXNlcl9pZCI6IjBmY2U5ZTc0LTAzMzgtNGZkMC1iMTI1LTVjYjE3NTBjN2FiZCIsIm9yaWdpbiI6InVhYSIsInVzZXJfbmFtZSI6ImFkbWluIiwiZW1haWwiOiJhZG1pbiIsImF1dGhfdGltZSI6MTQ4MzU0MjM0MSwicmV2X3NpZyI6IjJkZjU5MjhhIiwiaWF0IjoxNDgzNTQyMzQxLCJleHAiOjE0ODM2Mjg3NDEsImlzcyI6Imh0dHBzOi8vMTAuMjQ0LjY2LjI6ODQ0My9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJib3NoX2NsaSIsIm9wZW5pZCIsImJvc2giXX0.sZBoWKDwz--YP4u9VsXLCGW4pjbIRV1UD21grNj-VGemQGi26Yc-45CxXa_A3KWm52-pMBVwfhDt7tUxcHHtKNPkxa3a5MFmKkX0dzo-scWlw5uocNitviwelcTvOIzUnWUJMbOLuR0HlnGw_Vb3iXEUWTsSkTCMJiV0-oY8vp8IuUgJ9YchuZtSORh127dpJyV6ifKi9j9zPf9TbfmOIiFqS6Wp89xPY5zJTeXhPlK-ZL-qWIV6hb9aQqktGepFj31pb_p1qDD16IOxg-rlMxPxwNx_ov35cxfnrwRAHQt-Wlf4xVm2lGp5mTRDtyK1RoDItZz87GYPfdVbxPKwTleKhQq1fjMf_LN6-2OULA1QlIBxy0rPOSdythoMd3bTWJhrOvzijCgdI10-xFeevQlzZ2VdLHotH_kvYn5AIbYc_lefBc6liJoQXazTlzEyall7ZN2FConelh-oTbZxg1uwYY7SE2PBFSp_29OxjO5hA9ByeAQ2JAN6JFT0MNj7FCUEhyXYDyr1TM9BzgTsxkJBf0YzKZ1-58oqLglNnhb1Hv9BXfJ4jfM_OgMPE2Mp0BRyZffasZqWoqj1QQTS_raQc30b1r1476kJlyGC_itSyB_IWY13ORvooB8tqqQmWWPopo6IIrhSODn-YcVdjCskyyiMrq0j69KdLvVtESQ"))
			Expect(tokenResp.TokenType).Should(Equal("bearer"))
			Expect(tokenResp.ExpiresIn).Should(Equal(86399))
			Expect(tokenResp.Scope).Should(Equal("openid bosh.admin"))
			Expect(tokenResp.JTI).Should(Equal("c5a400529e01466d97a8187ed0c24e89"))
		})

	})

})
