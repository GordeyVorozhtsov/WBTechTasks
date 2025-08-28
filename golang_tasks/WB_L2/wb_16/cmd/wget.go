package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

var (
	site          string
	artifactsPath string = "artifacts/"
	maxDepth      int
	maxWorkers    int
	timeout       int
)

var wgetCmd = &cobra.Command{
	Use:   "wget",
	Short: "Website mirroring tool",
	Long:  "Download websites with all embedded content",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runWget,
}

func init() {
	rootCmd.RunE = runWget
	rootCmd.Flags().StringVarP(&site, "website", "w", "", "Website URL")
	rootCmd.Flags().IntVarP(&maxDepth, "depth", "d", 1, "Maximum recursion depth")
	rootCmd.Flags().IntVarP(&maxWorkers, "workers", "j", 1, "Maximum concurrent downloads")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 30, "Request timeout in seconds")
}

type DownloadManager struct {
	visited     map[string]bool // Посещенные URL
	visitedLock sync.RWMutex    // Блокировка для потокобезопасности
	baseURL     *url.URL
	baseDomain  string
	artifacts   string
	client      *http.Client
	queue       chan *DownloadTask // Очередь задач
	wg          sync.WaitGroup
}

// Сайт(задача) для загрузки
type DownloadTask struct {
	URL   string
	Depth int
}

func runWget(cmd *cobra.Command, args []string) error {
	fmt.Printf("Starting mirroring with depth %d and %d workers\n", maxDepth, maxWorkers)

	// Если URL указан через флаг -w, добавляем его в аргументы
	if site != "" {
		args = append([]string{site}, args...)
	}

	// Создаем менеджер загрузок
	manager, err := NewDownloadManager(args[0])
	if err != nil {
		return err
	}

	// Запускаем воркеров для обработки задач
	for i := 0; i < maxWorkers; i++ {
		go manager.worker()
	}

	// Добавляем начальные задачи из аргументов
	for _, u := range args {
		manager.addTask(&DownloadTask{
			URL:   normalizeURL(u),
			Depth: 0,
		})
	}

	go func() {
		manager.wg.Wait()
		close(manager.queue)
	}()

	manager.wg.Wait()

	// Сохраняем историю посещенных URL
	if err := manager.saveVisited(); err != nil {
		fmt.Printf("Warning: could not save visited history: %v\n", err)
	}

	fmt.Printf("Mirroring completed! Files saved to: %s\n", artifactsPath)
	return nil
}

// NewDownloadManager создает новый менеджер загрузок
func NewDownloadManager(startURL string) (*DownloadManager, error) {
	baseURL, err := url.Parse(normalizeURL(startURL))
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	dm := &DownloadManager{
		visited:    make(map[string]bool),
		baseURL:    baseURL,
		baseDomain: baseURL.Hostname(),
		artifacts:  artifactsPath,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
		queue: make(chan *DownloadTask, 1000),
	}

	// Загружаем историю посещенных URL если есть
	if err := dm.loadVisited(); err != nil {
		fmt.Printf("Warning: could not load visited history: %v\n", err)
	}
	return dm, nil
}

// addTask добавляет задачу в очередь если URL еще не посещен
func (dm *DownloadManager) addTask(task *DownloadTask) {
	dm.visitedLock.Lock()
	defer dm.visitedLock.Unlock()

	// Пропускаем если уже посещали или превышена глубина
	if dm.visited[task.URL] || task.Depth > maxDepth {
		return
	}
	// Помечаем как посещенный и добавляем в очередь
	dm.visited[task.URL] = true
	dm.wg.Add(1)
	dm.queue <- task
}

// worker обрабатывает задачи из очереди
func (dm *DownloadManager) worker() {
	for task := range dm.queue {
		dm.processTask(task)
		dm.wg.Done()
	}
}

// Обрабатываем одну задачу на загрузку
func (dm *DownloadManager) processTask(task *DownloadTask) {
	fmt.Printf("Processing: %s (depth: %d)\n", task.URL, task.Depth)

	// Загружаем ресурс
	content, contentType, err := dm.downloadResource(task.URL)
	if err != nil {
		fmt.Printf("Error downloading %s: %v\n", task.URL, err)
		return
	}

	// Сохраняем файл
	localPath, err := dm.saveFile(task.URL, content, contentType)
	if err != nil {
		fmt.Printf("Error saving %s: %v\n", task.URL, err)
		return
	}

	// Если это html и не превышена глубина -> парсим ссылки
	if strings.Contains(contentType, "text/html") && task.Depth < maxDepth {
		dm.parseAndQueueLinks(task.URL, content, task.Depth+1)
	}

	fmt.Printf("Saved: %s -> %s\n", task.URL, localPath)
}

// Загружаем ресурс по URL
func (dm *DownloadManager) downloadResource(urlStr string) ([]byte, string, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, "", err
	}

	// попытка установить заголовки, чтобы условный сайт wb не посчитал меня ботом и скачал все,
	// но чтото пошло не так и разбираться я с этим особо не стал
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")

	resp, err := dm.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status: %s", resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return content, resp.Header.Get("Content-Type"), nil
}

// Сохраняем содержимое в файл
func (dm *DownloadManager) saveFile(urlStr string, content []byte, contentType string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	localPath := dm.getLocalPath(parsedURL)

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return "", err
	}

	return localPath, nil
}

// Возвращаем локальный путь для сохранения файла
func (dm *DownloadManager) getLocalPath(parsedURL *url.URL) string {
	basePath := parsedURL.Path
	if basePath == "" || basePath == "/" {
		basePath = "/index.html"
	}

	// Обрабатываем query параметры в имени файла
	if parsedURL.RawQuery != "" {
		ext := filepath.Ext(basePath)
		if ext != "" {
			basePath = strings.TrimSuffix(basePath, ext) +
				"_" + strings.ReplaceAll(parsedURL.RawQuery, "&", "_") + ext
		} else {
			basePath += "_" + strings.ReplaceAll(parsedURL.RawQuery, "&", "_")
		}
	}

	// Сохраняем в структуре папок: artifacts/домен/путь
	return filepath.Join(dm.artifacts, parsedURL.Hostname(), basePath)
}

// парсит html и добавляет найденные ссылки в очередь
func (dm *DownloadManager) parseAndQueueLinks(baseURL string, content []byte, depth int) {
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "a", "link": // Ссылки и CSS
				dm.enqueueIfValid(baseURL, n, "href", depth)
			case "img", "script", "iframe": // Изображения, скрипты
				dm.enqueueIfValid(baseURL, n, "src", depth)
			case "form": // Формы
				dm.enqueueIfValid(baseURL, n, "action", depth)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinks(c)
		}
	}
	extractLinks(doc)
}

// Обрабатываем ссылку в очередь если она валидна
func (dm *DownloadManager) enqueueIfValid(baseURL string, n *html.Node, attr string, depth int) {
	for _, a := range n.Attr {
		if a.Key == attr {
			absoluteURL := dm.resolveURL(baseURL, a.Val)
			if dm.shouldDownload(absoluteURL) {
				dm.addTask(&DownloadTask{URL: absoluteURL, Depth: depth})
			}
			break
		}
	}
}

// Преобразуем относительный URL в абсолютный
func (dm *DownloadManager) resolveURL(baseURL, relativeURL string) string {
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return relativeURL
	}

	relative, err := url.Parse(relativeURL)
	if err != nil {
		return relativeURL
	}

	return base.ResolveReference(relative).String()
}

// Проверяем нужно ли скачивать URL
func (dm *DownloadManager) shouldDownload(urlStr string) bool {
	if strings.HasPrefix(urlStr, "mailto:") ||
		strings.HasPrefix(urlStr, "javascript:") ||
		strings.HasPrefix(urlStr, "tel:") ||
		strings.HasPrefix(urlStr, "data:") {
		return false
	}

	// Игнорируем якорные ссылки
	if strings.Contains(urlStr, "#") {
		return false
	}

	// Проверяем что ссылка с того же домена или относительная
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.Hostname() == dm.baseDomain || parsedURL.Hostname() == ""
}

// normalizeURL нормализует URL
func normalizeURL(urlStr string) string {
	urlStr = strings.TrimSpace(urlStr)

	// Добавляем https:// если нет протокола
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}

	// Убираем якоря
	if idx := strings.Index(urlStr, "#"); idx != -1 {
		urlStr = urlStr[:idx]
	}

	// Убираем trailing slash для единообразия
	if strings.HasSuffix(urlStr, "/") && len(urlStr) > 8 {
		urlStr = urlStr[:len(urlStr)-1]
	}
	return urlStr
}

// Сохраняем историю посещенных URL в файл
func (dm *DownloadManager) saveVisited() error {
	dm.visitedLock.RLock()
	defer dm.visitedLock.RUnlock()

	if len(dm.visited) == 0 {
		return nil
	}

	data, err := json.Marshal(dm.visited)
	if err != nil {
		return err
	}

	historyFile := filepath.Join(dm.artifacts, "visited_history.json")

	// Создаем директорию если нужно
	if err := os.MkdirAll(filepath.Dir(historyFile), 0755); err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

// Загружаем историю посещенных URL из файла
func (dm *DownloadManager) loadVisited() error {
	historyFile := filepath.Join(dm.artifacts, "visited_history.json")

	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(historyFile)
	if err != nil {
		return err
	}

	dm.visitedLock.Lock()
	defer dm.visitedLock.Unlock()
	return json.Unmarshal(data, &dm.visited)
}
