package helpers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func DoesStartWithFourDigits(postName string) bool {
	isNumeric := true
	for i := 0; i < 4; i++ {
		if !unicode.IsDigit(rune(postName[i])) {
			isNumeric = false
			break
		}
	}
	return isNumeric
}

func FileNameWithoutExt(fileName string) string {
	ext := filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(ext)]
}

func extractDriveFileID(url string) (string, error) {
	parts := strings.Split(url, "/")
	for i := range parts {
		if parts[i] == "d" && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}
	return "", fmt.Errorf("file ID not found")
}

func BuildGDriveImageUrl(shareUrl string) (string, error) {
	id, err := extractDriveFileID(shareUrl)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("https://drive.google.com/uc?export=download&id=%v", id), err
}

func DownloadImage(url string, filepath string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func CreateFolder(folderName string) error {
	err := os.Mkdir(folderName, 0777)
	if err != nil {
		return err
	}
	return nil
}

func NukeFolder(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return os.MkdirAll(path, 0o755)
}

func SetToken(tokenEnvName, tokenEnvPath string) (string, error) {
	token := os.Getenv(tokenEnvName)
	if token == "" {
		tokenPath := os.Getenv(tokenEnvPath)
		if tokenPath == "" {
			return "", errors.New(fmt.Sprintf("%v or %v environment variable are not set", tokenEnvName, tokenEnvPath))
		}
		tokenBytes, err := os.ReadFile(tokenPath)
		if err != nil {
			return "", errors.New(fmt.Sprintf("could not read git token from file %v: %v", tokenPath, err))
		}
		token = string(tokenBytes)
	}
	return token, nil
}
