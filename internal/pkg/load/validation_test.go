package load

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestValidation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Load")
}

var _ = Describe("Resource Validation", func() {
	It("Validates Resources", func() {
		_, err := PerformResourceValidation("/home/sjakati/go/src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/resources")
		Expect(err).ToNot(HaveOccurred())
	})
})
