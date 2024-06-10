package main

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGetFilePath тестирует функцию getFilePath
func TestGetFilePath(t *testing.T) {
	u, err := url.Parse("http://example.com/test")
	if err != nil {
		t.Fatalf("Failed to parse URL: %v", err)
	}
	outputDir := "output"

	expected := filepath.Join(outputDir, "example.com", "test", "index.html")
	actual := getFilePath(u, outputDir)

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

// TestResolveURL тестирует функцию resolveURL
func TestResolveURL(t *testing.T) {
	base := "http://example.com/test/"
	href := "page.html"

	expected := "http://example.com/test/page.html"
	actual := resolveURL(base, href)

	if actual != expected {
		t.Errorf("Expected %s but got %s", expected, actual)
	}
}

// TestSavePage тестирует функцию savePage
func TestSavePage(t *testing.T) {
	u := "http://example.com/test"
	outputDir := "output_test"
	os.RemoveAll(outputDir) // Удаляем директорию перед началом теста

	// Создание документа goquery
	htmlContent := `<html><head><title>Test</title></head><body><p>Hello, World!</p></body></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to create document: %v", err)
	}

	err = savePage(u, outputDir, doc)
	if err != nil {
		t.Fatalf("Failed to save page: %v", err)
	}

	expectedFilePath := filepath.Join(outputDir, "example.com", "test", "index.html")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected file %s does not exist", expectedFilePath)
	}

	// Проверка содержимого файла
	content, err := ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if !strings.Contains(string(content), "Hello, World!") {
		t.Errorf("Expected file content to contain 'Hello, World!' but got %s", string(content))
	}

	// Удаляем директорию после теста
	os.RemoveAll(outputDir)
}

// Mock HTTP Client
type mockClient struct {
	response *http.Response
	err      error
	c        *http.Client
}

func (m *mockClient) Get(url string) (*http.Response, error) {
	return m.response, m.err
}

// TestDownloadPage тестирует функцию downloadPage
func TestDownloadPage(t *testing.T) {
	htmlContent := `<html><head><title>Test</title></head><body><p>Hello, World!</p><a href="/page2.html">Page 2</a></body></html>`
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(strings.NewReader(htmlContent)),
	}

	client := &mockClient{response: resp}
	outputDir := "output_test"
	visited := make(map[string]bool)

	// Удаляем директорию перед началом теста
	os.RemoveAll(outputDir)

	// Тестовая страница
	downloadPage(client.c, "http://example.com/test", outputDir, visited)

	expectedFilePath := filepath.Join(outputDir, "example.com", "test", "index.html")
	if _, err := os.Stat(expectedFilePath); os.IsNotExist(err) {
		t.Fatalf("Expected file %s does not exist", expectedFilePath)
	}

	// Проверка содержимого файла
	content, err := ioutil.ReadFile(expectedFilePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if !strings.Contains(string(content), "Hello, World!") {
		t.Errorf("Expected file content to contain 'Hello, World!' but got %s", string(content))
	}

	// Удаляем директорию после теста
	os.RemoveAll(outputDir)
}
