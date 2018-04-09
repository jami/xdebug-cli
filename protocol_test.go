package xdebugcli_test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/jami/xdebug-cli"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func loadProtocolAsset(name string) []byte {
	b, err := ioutil.ReadFile(path.Join("_fixture", "protocol", name))
	if err != nil {
		Expect(err).Should(BeNil())
	}

	return b
}

func TestProtocol(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Protocol xml parsing test suite")
}

//
var _ = Describe("Protocol xml parsing", func() {

	Context("Test init protocol", func() {
		It("test init.1.xml", func() {
			assetData := loadProtocolAsset("init.1.xml")
			proto, err := xdebugcli.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&xdebugcli.ProtocolInit{}))

			initProto := proto.(*xdebugcli.ProtocolInit)
			Ω(initProto.FileURI).Should(Equal("file:///PhpProject1/index.php"))
			Ω(initProto.Language).Should(Equal("PHP"))
			Ω(initProto.AppID).Should(Equal("24001"))
			Ω(initProto.IDEKey).Should(Equal("jami"))
		})
	})

	Context("Test response protocol", func() {
		It("test response.bplist.1.xml", func() {
			assetData := loadProtocolAsset("response.bplist.1.xml")
			proto, err := xdebugcli.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&xdebugcli.ProtocolResponse{}))
			/*
				initProto := proto.(*xdebugcli.ProtocolInit)
				Ω(initProto.FileURI).Should(Equal("file:///PhpProject1/index.php"))
				Ω(initProto.Language).Should(Equal("PHP"))
				Ω(initProto.AppID).Should(Equal("24001"))
				Ω(initProto.IDEKey).Should(Equal("jami")) */
		})
	})
})
