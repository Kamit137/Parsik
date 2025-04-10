package ali

import (
	"encoding/json"
	"fmt"
	"log"
	// "net/url"
	"strings"

	"github.com/gocolly/colly"
	"github.com/playwright-community/playwright-go"
)

// JSON - структура для aer_data
type JSON struct {
	Widgets []widget `json:"widgets"`
}

type widget struct {
	WidgetID string   `json:"widgetId"`
	Children []widget `json:"children"`
	Props    props    `json:"props"`
}

type props struct {
	Name    string  `json:"name"`
	SkuInfo skuInfo `json:"skuInfo"`
}

type skuInfo struct {
	PropertyList []propertyList `json:"propertyList"`
	PriceList    []priceList    `json:"priceList"`
}

type propertyList struct {
	Name   string          `json:"name"`
	Values []propertyValue `json:"values"`
	ID     string          `json:"id"`
}

type propertyValue struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type priceList struct {
	ActivityAmount activityAmount `json:"activityAmount"`
	SkuId          string         `json:"skuId"`
}

type activityAmount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

func fail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func M() {
	link := "https://aliexpress.ru/item/1005008579400145.html?sku_id=12000045808902946&spm=a2g2w.productlist.search_results.3.756e396aH5OdPx"
	skuID := "12000045808902946"

	// Используем Playwright для рендеринга страницы
	pw, err := playwright.Run()
	fail(err, "playwright run failed")
	browser, err := pw.Chromium.Launch()
	fail(err, "browser launch failed")
	page, err := browser.NewPage()
	fail(err, "new page failed")
	_, err = page.Goto(link)
	fail(err, "page goto failed")

	// Получаем HTML страницы после рендеринга
	// html, err := page.Content()
	fail(err, "get content failed")

	// Закрываем браузер
	err = browser.Close()
	fail(err, "browser close failed")
	err = pw.Stop()
	fail(err, "playwright stop failed")

	// Парсим HTML с помощью Colly
	c := colly.NewCollector()

	c.OnHTML("#__AER_DATA__", func(e *colly.HTMLElement) {
		var result JSON
		err := json.Unmarshal([]byte(e.Text), &result)
		fail(err, "JSON unmarshal failed")

		for i := range result.Widgets {
			findPrice(result.Widgets[i], skuID)
		}
	})

	// Обработчик ошибок
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\nError: %v", r.Request.URL, r, err)
	})

	// Запускаем парсинг HTML, полученного от Playwright
	c.Visit(link)

	// Добавляем ожидание завершения всех горутин
	c.Wait()
}

func findPrice(w widget, skuID string) {
	if strings.Contains(w.WidgetID, "SnowProductContextWidget") {
		for _, price := range w.Props.SkuInfo.PriceList {
			if price.SkuId == skuID {
				fmt.Printf("Product: %s, Price: %.2f %s\n", w.Props.Name, price.ActivityAmount.Value, price.ActivityAmount.Currency)
				return
			}
		}
	}

	// Рекурсивно проверяем дочерние элементы
	for _, child := range w.Children {
		findPrice(child, skuID)
	}
}
