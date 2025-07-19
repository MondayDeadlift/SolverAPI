package codewars

import (
	"SolverAPI/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseURL     string
	httpClient  *http.Client
	kataBuffer  []string   // Буфер ID задач
	lastUpdated time.Time  // Время последнего обновления
	bufferMutex sync.Mutex // Для потокобезопасности
}

func NewClient(baseURL string) *Client {
	c := &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Первоначальное заполнение буфера
	go c.RefreshBuffer()
	return c
}

// Автоматическое обновление буфера
func (c *Client) RefreshBuffer() {
	c.bufferMutex.Lock()
	defer c.bufferMutex.Unlock()

	if time.Since(c.lastUpdated) < 1*time.Hour && len(c.kataBuffer) > 0 {
		return // Не обновляем чаще чем раз в час
	}

	ids, err := c.scrapeKataList()
	if err != nil {
		log.Printf("Failed to refresh kata buffer: %v", err)
		return
	}

	c.kataBuffer = ids
	c.lastUpdated = time.Now()
	log.Printf("Kata buffer refreshed, %d tasks available", len(ids))
}

func (c *Client) GetUser(username string) (*model.CodewarsUser, error) {
	url := fmt.Sprintf("%s/users/%s", c.baseURL, username)
	log.Printf("Requesting user from URL: %s", url) // Добавьте это

	//Создаем GET-запрос с контекстом
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	//Выполняем запрос через наш httpClient (с таймаутом)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close() // Важно закрывать тело ответа

	//Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//Читаем и парсим JSON
	var user model.CodewarsUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user %q not found on Codewars", username)
	}

	// 5. Возвращаем результат
	return &user, nil
}

// GetKata возвращает информацию о задаче по ID
func (c *Client) GetKata(ctx context.Context, id string) (*model.CodewarsKata, error) {
	url := fmt.Sprintf("%s/code-challenges/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var kata model.CodewarsKata
	if err := json.NewDecoder(resp.Body).Decode(&kata); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &kata, nil
}

// Метод для получения конкретной задачи
func (c *Client) GetKataByID(ctx context.Context, id string) (*model.CodewarsKata, error) {
	url := fmt.Sprintf("%s/code-challenges/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var kata model.CodewarsKata
	if err := json.NewDecoder(resp.Body).Decode(&kata); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &kata, nil
}

// метод для заполнения буфера
func (c *Client) fillKataBuffer(ctx context.Context) error {
	c.bufferMutex.Lock()
	defer c.bufferMutex.Unlock()

	url := fmt.Sprintf("%s/code-challenges?page=0&pageSize=50", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	c.kataBuffer = make([]string, 0, len(data.Data))
	for _, kata := range data.Data {
		c.kataBuffer = append(c.kataBuffer, kata.ID)
	}

	return nil
}

// Новый метод для получения случайной задачи
func (c *Client) GetRandomKataFromBuffer(ctx context.Context) (*model.CodewarsKata, error) {
	id, err := c.GetRandomKataID(ctx)
	if err != nil {
		return nil, err
	}
	return c.GetKataByID(ctx, id)
}

func (c *Client) scrapeKataList() ([]string, error) {
	//Создаем HTTP-запрос с таймаутом
	req, err := http.NewRequest("GET", "https://www.codewars.com/kata/search", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки, чтобы имитировать браузер
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html")

	// Выполняем запрос с таймаутом
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("scrape request failed: %w", err)
	}
	defer resp.Body.Close()

	//Проверяем статус код
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	// Читаем тело ответа с ограничением по размеру
	maxBytes := 10 * 1024 * 1024 // 10MB
	body, err := io.ReadAll(io.LimitReader(resp.Body, int64(maxBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Регулярка для поиска ID задач в HTML
	re := regexp.MustCompile(`/kata/([a-f0-9]{24})`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	if len(matches) == 0 {
		return nil, errors.New("no kata IDs found in HTML")
	}

	// Убираем дубликаты через map
	uniqueIDs := make(map[string]struct{})
	for _, match := range matches {
		if len(match) > 1 {
			uniqueIDs[match[1]] = struct{}{}
		}
	}

	//Конвертируем в слайс
	result := make([]string, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		result = append(result, id)
	}

	// Проверяем, что нашли хотя бы несколько задач
	if len(result) < 5 {
		return nil, fmt.Errorf("found too few katas (%d), possible parsing error", len(result))
	}

	return result, nil
}

func (c *Client) GetRandomKataID(ctx context.Context) (string, error) {
	c.bufferMutex.Lock()
	defer c.bufferMutex.Unlock()

	if len(c.kataBuffer) == 0 {
		if err := c.fillKataBuffer(ctx); err != nil {
			return "", fmt.Errorf("failed to fill buffer: %w", err)
		}
	}

	return c.kataBuffer[rand.Intn(len(c.kataBuffer))], nil
}
