package dbgp_test

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/jami/xdebug-cli/dbgp"
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
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolInit{}))

			initProto := proto.(*dbgp.ProtocolInit)
			Ω(initProto.FileURI).Should(Equal("file:///PhpProject1/index.php"))
			Ω(initProto.Language).Should(Equal("PHP"))
			Ω(initProto.AppID).Should(Equal("24001"))
			Ω(initProto.IDEKey).Should(Equal("jami"))
		})
	})

	Context("Test response protocol", func() {
		It("test response.bplist.1.xml", func() {
			assetData := loadProtocolAsset("response.bplist.1.xml")
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolResponse{}))
			responseProto := proto.(*dbgp.ProtocolResponse)

			Ω(responseProto.Command).Should(Equal("breakpoint_list"))
			Ω(responseProto.BreakpointList).Should(HaveLen(0))
		})

		It("test response.bplist.2.xml", func() {
			assetData := loadProtocolAsset("response.bplist.2.xml")
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolResponse{}))
			responseProto := proto.(*dbgp.ProtocolResponse)

			Ω(responseProto.Command).Should(Equal("breakpoint_list"))
			Ω(responseProto.BreakpointList).Should(HaveLen(1))

			Ω(responseProto.BreakpointList).Should(Equal([]dbgp.ProtocolBreakpoint{{
				Type:     "line",
				Line:     19,
				State:    "enabled",
				FileName: "file:///home/jami/NetBeansProjects/PhpProject1/index.php",
				HitCount: 0,
			}}))
		})

		It("test response.bpset.1.xml", func() {
			assetData := loadProtocolAsset("response.bpset.1.xml")
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolResponse{}))
			responseProto := proto.(*dbgp.ProtocolResponse)

			Ω(responseProto.Command).Should(Equal("breakpoint_set"))
		})

		It("test response.run.1.xml", func() {
			assetData := loadProtocolAsset("response.run.1.xml")
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolResponse{}))
			responseProto := proto.(*dbgp.ProtocolResponse)

			Ω(responseProto.Command).Should(Equal("run"))
			Ω(responseProto.Status).Should(Equal("break"))
			Ω(responseProto.Reason).Should(Equal("ok"))
			Ω(responseProto.Message.Filename).Should(Equal("file:///PhpProject1/index.php"))
			Ω(responseProto.Message.Line).Should(Equal(19))
		})

		It("test response.sget.1.xml", func() {
			assetData := loadProtocolAsset("response.sget.1.xml")
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolResponse{}))
			responseProto := proto.(*dbgp.ProtocolResponse)

			Ω(responseProto.Command).Should(Equal("stack_get"))
			Ω(responseProto.StackList).Should(HaveLen(1))

			stack := responseProto.StackList[0]
			Ω(stack.Filename).Should(Equal("file:///PhpProject1/index.php"))
			Ω(stack.Level).Should(Equal(0))
			Ω(stack.Line).Should(Equal(19))
			Ω(stack.Type).Should(Equal("file"))
			Ω(stack.Where).Should(Equal("{main}"))
		})

		It("test response.cn.1.xml", func() {
			assetData := loadProtocolAsset("response.cn.1.xml")
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolResponse{}))
			responseProto := proto.(*dbgp.ProtocolResponse)

			Ω(responseProto.Command).Should(Equal("context_names"))
			Ω(responseProto.ContextList).Should(HaveLen(3))

			Ω(responseProto.ContextList[0]).Should(Equal(dbgp.ProtocolContext{
				ID:   "0",
				Name: "Locals",
			}))

			Ω(responseProto.ContextList[1]).Should(Equal(dbgp.ProtocolContext{
				ID:   "1",
				Name: "Superglobals",
			}))

			Ω(responseProto.ContextList[2]).Should(Equal(dbgp.ProtocolContext{
				ID:   "2",
				Name: "User defined constants",
			}))
		})

		It("test response.cg.1.xml", func() {
			assetData := loadProtocolAsset("response.cg.1.xml")
			proto, err := dbgp.CreateProtocolFromXML(string(assetData))
			Ω(err).ShouldNot(HaveOccurred())
			Ω(proto).Should(BeAssignableToTypeOf(&dbgp.ProtocolResponse{}))
			responseProto := proto.(*dbgp.ProtocolResponse)

			Ω(responseProto.Command).Should(Equal("context_get"))
			Ω(responseProto.Context).Should(Equal("1"))
			Ω(responseProto.PropertyList).Should(HaveLen(7))

			Ω(responseProto.PropertyList[0]).Should(Equal(dbgp.ProtocolProperty{
				Name:        "$_COOKIE",
				Fullname:    "$_COOKIE",
				Type:        "array",
				Children:    0,
				NumChildren: 0,
				Page:        0,
				PageSize:    32,
				Content:     "",
				Property:    nil,
			}))
		})
	})
})
