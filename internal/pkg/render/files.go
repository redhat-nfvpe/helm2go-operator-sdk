package render

import (
	"os"
	"path/filepath"
)

// takes a map of injected files and writes them to a temp directory
func writeToTemp(files map[string]string, outParentDir string) (string, error) {
	os.RemoveAll(filepath.Join(outParentDir, "temp"))
	err := os.Mkdir(filepath.Join(outParentDir, "temp"), 0700)
	if err != nil {
		return "", err
	}
	for filePath, fileContent := range files {
		// open the file
		// must create the path if it does not exist
		d := filepath.Dir(filepath.Join(outParentDir, "temp", filePath))
		if _, err := os.Stat(d); os.IsNotExist(err) {
			os.MkdirAll(d, 0700)
		}
		f, err := os.Create(filepath.Join(outParentDir, "temp", filePath))
		if err != nil {
			return "", err
		}
		// write the file
		_, err = f.WriteString(fileContent)
		if err != nil {
			return "", err
		}
		f.Close()
	}
	// return the temp directory filepath
	return filepath.Join(outParentDir, "temp"), err
}
