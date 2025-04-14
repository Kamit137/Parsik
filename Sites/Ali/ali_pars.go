package ali

import (
	"fmt"
	// "log"
	"strconv"
	"strings"
	// "time"

	"github.com/gocolly/colly"
	"github.com/playwright-community/playwright-go"
)

func getRenderedHTML(url string) (string, error) {
	pw, err := playwright.Run()
	if err != nil {
		return "", fmt.Errorf("playwright launch error: %v", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("browser launch error: %v", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return "", fmt.Errorf("new page error: %v", err)
	}

	// Устанавливаем язык и регион
	page.SetExtraHTTPHeaders(map[string]string{
		"Accept-Language": "ru-RU,ru;q=0.9",
	})

	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000),
	})
	if err != nil {
		return "", fmt.Errorf("navigation error: %v", err)
	}

	return page.Content()
}

func ParseProduct(url string) (*JSON, error) {
	html, err := getRenderedHTML(url)
	if err != nil {
		return nil, fmt.Errorf("render error: %v", err)
	}

	c := colly.NewCollector()
	result := &JSON{Widgets: make([]Widget, 0)}

	// Обновленные селекторы для AliExpress 2025
	c.OnHTML("h1[itemprop='name']", func(e *colly.HTMLElement) {
		result.Widgets = append(result.Widgets, Widget{
			WidgetID: "title",
			Props: Props{
				Name: strings.TrimSpace(e.Text),
			},
		})
	})

	c.OnHTML("div[data-pl='product-price']", func(e *colly.HTMLElement) {
		priceStr := strings.NewReplacer("₽", "", ",", "", " ", "").Replace(e.Text)
		if price, err := strconv.Atoi(priceStr); err == nil {
			result.Widgets = append(result.Widgets, Widget{
				WidgetID: "price",
				Props: Props{
					SkuInfo: SkuInfo{
						ActivityAmount: ActivityAmount{Value: price},
					},
				},
			})
		}
	})

	c.OnHTML("img[itemprop='image']", func(e *colly.HTMLElement) {
		if img := e.Attr("src"); img != "" {
			result.Widgets = append(result.Widgets, Widget{
				WidgetID: "image",
				Props: Props{
					SkuInfo: SkuInfo{
						SkuId: "https:" + strings.Split(img, "?")[0],
					},
				},
			})
		}
	})

	c.OnResponse(func(r *colly.Response) {
		r.Body = []byte(html)
	})

	if err := c.Visit("http://localhost/render"); err != nil {
		return nil, fmt.Errorf("parsing error: %v", err)
	}

	if len(result.Widgets) == 0 {
		return nil, fmt.Errorf("no data found after rendering")
	}

	return result, nil
}
