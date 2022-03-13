package rwFile

import (
	"os"
)

func WriteFile(content, path string) error {
	File, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	_, err := File.WriteString(content)
	return err
}

func ReadFile(path string) string {
	content, _ := os.ReadFile(path)
	return string(content)
}
