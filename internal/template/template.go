package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

//go:embed templates/*
var fs embed.FS

const templatesPath = "templates"

type Details struct {
	Date  time.Time
	Venue Venue
}

type Venue struct {
	Name    string
	Address string
	City    string
}

type Post struct {
	ID           string
	Name         string
	Description  string
	Author       string
	Ensemble     string
	Tags         []string
	Categories   []string
	ImageName    string
	BuyTicketUrl string
	FirstDate    time.Time
	Details      []*Details
}

func GeneratePost(templateFileName, folder, postFileName string, post Post) error {
	tmpl, err := createTemplate(templateFileName)
	if err != nil {
		return err
	}
	f, err := createFile(folder, postFileName)
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(f, "post", post)
	if err != nil {
		return err
	}
	return nil
}

func createTemplate(templateName string) (*template.Template, error) {
	path := filepath.Join(templatesPath, templateName)
	tmpl, err := template.New("post").ParseFS(fs, path)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func createFile(postFolder, postName string) (*os.File, error) {
	fileName := fmt.Sprintf("%s", postName)
	filePath := filepath.Join(postFolder, fileName)
	f, err := os.Create(filePath)
	return f, err
}
