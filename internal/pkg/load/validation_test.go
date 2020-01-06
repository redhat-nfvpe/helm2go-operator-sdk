package load

import (
	"os"
	"path/filepath"
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
		gopath := os.Getenv("GOPATH")
		_, err := PerformResourceValidation(filepath.Join(gopath, "/src/github.com/redhat-nfvpe/helm2go-operator-sdk/test/resources"),false)
		Expect(err).ToNot(HaveOccurred())
	})
})
