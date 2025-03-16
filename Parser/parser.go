package pars_wb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// JSON представляет структуру JSON, которую мы ожидаем получить.
type JSON struct {
	Data Data `json:"data"`
}

type Data struct {
	Products []Product `json:"products"`
}

type Product struct {
	Name       string `json:"name"`
	SalePriceU int    `json:"salePriceU"`
}

func fail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func Wb(productURL string)(JSON) {
	url, err := extractWildberriesCardURL(productURL)
	result, err := http.Get(url)
	fail(err, "request EB fail")

	if result.StatusCode > 299 {
		fmt.Println("error ", result.StatusCode)
		os.Exit(0)
	}

	body, err := io.ReadAll(result.Body)
	fail(err, "read body")

	fmt.Printf("%s\n", body)

	var jsonData JSON

	err = json.Unmarshal(body, &jsonData)
	fail(err, "json parser")

	return jsonData
}

func extractWildberriesCardURL(productURL string) (string, error) {
	parsedURL, err := url.Parse(productURL)
	if err != nil {
		return "", fmt.Errorf("некорректный URL: %w", err)
	}

	// Извлекаем ID товара (nm) из параметров запроса
	queryValues := parsedURL.Query()
	nm := queryValues.Get("nm")

	if nm == "" {
		// Попробуем извлечь из пути, если нет в параметрах
		pathParts := strings.Split(parsedURL.Path, "/")
		if len(pathParts) > 2 {
			lastPart := pathParts[len(pathParts)-2] // Use second to last element
			// Проверяем, что последняя часть - это число (nm)

			nm = lastPart

		} else {
			return "", fmt.Errorf("не удалось извлечь ID товара (nm)")
		}
	}

	// Формируем ссылку на карточку товара
	cardURL := fmt.Sprintf("https://card.wb.ru/cards/detail?appType=1&curr=rub&dest=-1257786&spp=27&nm=%s", nm)

	return cardURL, nil
}
