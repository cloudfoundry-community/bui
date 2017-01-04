package api_test

import (
	"io/ioutil"
	"log"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v2"

	. "github.com/cloudfoundry-community/bui/api"
)

var _ = Describe("API Configuration", func() {
	Describe("Configuration", func() {
		var a *Api

		BeforeEach(func() {
			a = NewApi()
			Ω(a).ShouldNot(BeNil())
			setup(MockRoute{"GET", "/info", info, ""})
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
			config := Config{
				BoshAddr: server.URL,
			}
			configByte, err := yaml.Marshal(config)
			Expect(err).Should(BeNil())
			tmpfile, err := ioutil.TempFile("", "test")
			if err != nil {
				log.Fatal(err)
			}

			defer os.Remove(tmpfile.Name()) // clean up

			if _, err := tmpfile.Write(configByte); err != nil {
				log.Fatal(err)
			}
			if err := tmpfile.Close(); err != nil {
				log.Fatal(err)
			}
			Ω(a.ReadConfig(tmpfile.Name())).Should(Succeed())
			Ω(a.Web.Addr).Should(Equal(":9304"))
		})

	})
})
