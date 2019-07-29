package render

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRender(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Render")
}

var _ = Describe("Write To Temp", func() {
	It("Writes Correctly", func() {
		testMap := map[string]string{
			"oneDir/test.txt":        "This is the content within the file",
			"twoDir/src/another.txt": "This is more file content",
		}

		testDir, err := filepath.Abs("../../../test/")
		Expect(err).ToNot(HaveOccurred())

		_, err = InjectedToTemp(testMap, testDir)
		Expect(err).ToNot(HaveOccurred())

		_, err = os.Stat(filepath.Join(testDir, "temp"))
		Expect(os.IsNotExist(err)).To(BeFalse())

		_, err = os.Stat(filepath.Join(testDir, "temp", "oneDir", "test.txt"))
		Expect(os.IsNotExist(err)).To(BeFalse())

		_, err = os.Stat(filepath.Join(testDir, "temp", "twoDir", "src"))
		Expect(os.IsNotExist(err)).To(BeFalse())

		err = os.RemoveAll(filepath.Join(testDir, "temp"))
		if err != nil {
			panic(err)
		}
	})
})
