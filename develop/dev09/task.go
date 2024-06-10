package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

/*
Утилита wget

Реализовать утилиту wget с возможностью скачивать сайты целиком.

*/

func savePage(u string, outputDir string, doc *goquery.Document) error {
	// парсинг URL
	parseURL, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("failed to parse URL %s: %w", u, err)
	}

	// filePath -  Построение пути для сохранения файла
	filePath := getFilePath(parseURL, outputDir)

	// Создание всех необходимых директорий для сохранения файла
	err = os.MkdirAll(path.Dir(filePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories for %s: %w", filePath, err)
	}

	// Создание файла для записи HTML содержимого
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}

	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("failed to close file %s: %v", filePath, err)
		}
	}()

	// Получение HTML содержимого из goquery.Document
	html, err := doc.Html()
	if err != nil {
		return fmt.Errorf("failed to get HTML conntent for %s: %w", u, err)
	}

	// Запись HTML содержимого в файл
	_, err = file.WriteString(html)
	if err != nil {
		return fmt.Errorf("failed to write HTML conntent for %s: %w", filePath, err)
	}

	fmt.Println("Saved:", filePath)
	return nil
}

// getFilePath конструирует путь файла из URL и директории для сохранения
func getFilePath(u *url.URL, outputDir string) string {
	filePath := path.Join(outputDir, u.Host, u.Path)
	// Если путь не содержит расширения, добавляем "index.html"
	if path.Ext(filePath) == "" {
		filePath = path.Join(filePath, "index.html")
	}

	return filePath

}

// downloadPage загружает и обрабатывает одну страницу
func downloadPage(client *http.Client, u string, outputDir string, visited map[string]bool) {
	// проверка, был ли уже посещен этот URL

	if visited[u] {
		return
	}

	visited[u] = true

	// Выполнение HTTP запроса
	resp, err := client.Get(u)
	if err != nil {
		log.Printf("error downloading page %s: %v", u, err)
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			log.Printf("error closing response body for %s: %v", u, err)
		}
	}()

	// проверка статуса ответа
	if resp.StatusCode != 200 {
		log.Printf("non-OK HTTP status %d for %s", resp.StatusCode, u)
	}

	// Создание документа из тела ответа
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("error parsing page %s: %v", u, err)
	}

	// создание страницы
	if err = savePage(u, outputDir, doc); err != nil {
		log.Printf("error saving page %s: %v", u, err)
		return
	}

	//поиск и обработка ссылок

	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if !exists {
			return
		}

		// Разрешение абсолютного URL из относительного
		absURL := resolveURL(u, link)
		if absURL != "" && strings.HasPrefix(absURL, u) {
			downloadPage(client, absURL, outputDir, visited) // Рекурсивная загрузка ссылок
		}

	})
}

// // resolveURL разрешает потенциально относительный URL относительно базового UR
func resolveURL(base, href string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		log.Printf("error parsing base URL %s: %v", base, err)
		return ""
	}

	hrefURL, err := url.Parse(href)
	if err != nil {
		log.Printf("error parsing href URL %s: %v", href, err)
		return ""
	}

	return baseURL.ResolveReference(hrefURL).String()
}

func main() {

	startURL := flag.String("url", "", "URL of the site to download")
	outputDir := flag.String("output", "output", "Directory to save files")
	flag.Parse()

	if *startURL == "" {
		fmt.Println("Usage: wget -url <url> -output <output>")
		return
	}

	err := os.MkdirAll(*outputDir, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	client := &http.Client{}
	visited := make(map[string]bool)
	downloadPage(client, *startURL, *outputDir, visited)
}
