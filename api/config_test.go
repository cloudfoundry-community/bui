package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry-community/bui/api"
)

var _ = Describe("API Configuration", func() {
	Describe("Configuration", func() {
		var a *Api

		BeforeEach(func() {
			a = NewApi()
			Ω(a).ShouldNot(BeNil())
		})

		It("handles missing files", func() {
			Ω(a.ReadConfig("/path/to/nowhere")).ShouldNot(Succeed())
		})

		It("handles malformed YAML files", func() {
			Ω(a.ReadConfig("test/etc/config.xml")).ShouldNot(Succeed())
		})

		It("handles YAML files with missing directives", func() {
			Ω(a.ReadConfig("test/etc/empty.yml")).Should(Succeed())
			Ω(a.Web.Addr).Should(Equal(":9304"))
		})

		It("handles YAML files with all the directives", func() {
			Ω(a.ReadConfig("test/etc/valid.yml")).Should(Succeed())
			Ω(a.Web.Addr).Should(Equal(":8988"))
		})

	})
})
