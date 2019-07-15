package render

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteToTemp(t *testing.T) {

	testMap := map[string]string{
		"oneDir/test.txt":        "This is the content within the file",
		"twoDir/src/another.txt": "This is more file content",
	}

	testDir, err := filepath.Abs("../../test/")
	if err != nil {
		panic(err)
	}

	_, err = InjectedToTemp(testMap, testDir)
	if err != nil {
		t.Errorf("Unexpected Error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(testDir, "temp")); os.IsNotExist(err) {
		t.Errorf("Directory %v Does Not Exist!", filepath.Join(testDir, "temp"))
	}

	if _, err := os.Stat(filepath.Join(testDir, "temp", "oneDir", "test.txt")); os.IsNotExist(err) {
		t.Errorf("Directory %v Does Not Exist!", filepath.Join(testDir, "temp", "oneDir", "test.txt"))
	}

	if _, err := os.Stat(filepath.Join(testDir, "temp", "twoDir", "src")); os.IsNotExist(err) {
		t.Errorf("Directory %v Does Not Exist!", "twoDir")
	}

	// clean up
	if err = os.RemoveAll(filepath.Join(testDir, "temp")); err != nil {
		t.Error(err)
	}
}
