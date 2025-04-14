package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

func main() {
	// 1. Инициализация Playwright с игнорированием ошибок зависимостей
	pw, err := playwright.Run(&playwright.RunOptions{
		DriverDirectory: "/tmp",
	})
	if err != nil {
		log.Fatalf("Ошибка запуска Playwright: %v", err)
	}
	defer pw.Stop()

	// 2. Настройка Firefox с обходом защиты
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
		Args: []string{
			"--disable-blink-features=AutomationControlled",
			"--profile", "/tmp/firefox_profile",
		},
		FirefoxUserPrefs: map[string]interface{}{
			"general.useragent.override":   generateUserAgent(),
			"privacy.resistFingerprinting": false,
			"media.peerconnection.enabled": false,
			"webgl.disabled":               true,
		},
	})
	if err != nil {
		log.Fatalf("Ошибка запуска Firefox: %v", err)
	}
	defer browser.Close()

	// 3. Создание контекста с резидентным прокси (ОБЯЗАТЕЛЬНО раскомментируйте!)
	/*
		context, err := browser.NewContext(playwright.BrowserNewContextOptions{
			Proxy: &playwright.Proxy{
				Server:   "http://ваш_резидентный_прокси:порт",
				Username: "логин",
				Password: "пароль",
			},
		})
	*/
	context, err := browser.NewContext()
	if err != nil {
		log.Fatalf("Ошибка создания контекста: %v", err)
	}
	defer context.Close()

	// 4. Создание страницы с увеличенными таймаутами
	page, err := context.NewPage()
	if err != nil {
		log.Fatalf("Ошибка создания страницы: %v", err)
	}

	// 5. Улучшенная эмуляция человеческого поведения
	simulateHumanBehavior(page)

	// 6. Переход на страницу с рандомизированными параметрами
	_, err = page.Goto("https://www.ozon.ru", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(180000), // 3 минуты
		Referer:   playwright.String("https://www.google.com/"),
	})
	if err != nil {
		log.Fatalf("Ошибка загрузки страницы: %v", err)
	}

	// 7. Проверка на блокировку
	if isBlocked(page) {
		log.Println("ВАЖНО: Для успешного парсинга необходимо:")
		log.Println("1. Раскомментировать и настроить резидентный прокси в коде")
		log.Println("2. Установить Headless: true для production")
		log.Println("3. Увеличить задержки в simulateHumanBehavior()")
		return
	}

	// 8. Парсинг данных
	title := extractWithFallback(page, []string{
		"h1[data-widget='webProductHeading']",
		"h1.b0v",
		"h1.m8m_28",
		"h1",
	}, "Название")

	price := extractWithFallback(page, []string{
		"div[data-widget='webPrice']",
		"span.l3z_28",
		"div[data-widget='price']",
		"[class*='Price']",
	}, "Цена")

	// 9. Вывод результатов
	fmt.Println("\n=== Успешный парсинг ===")
	fmt.Printf("Название: %s\n", strings.TrimSpace(title))
	fmt.Printf("Цена: %s\n", strings.TrimSpace(price))
}

func generateUserAgent() string {
	agents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_5_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	return agents[rand.Intn(len(agents))]
}

func simulateHumanBehavior(page playwright.Page) {
	rand.Seed(time.Now().UnixNano())

	// Длинные случайные задержки
	time.Sleep(time.Duration(rand.Intn(20)+10) * time.Second)

	// Разнообразные действия
	actions := []func(){
		func() { page.Mouse().Move(float64(rand.Intn(1200)), float64(rand.Intn(800))) },
		func() { page.Keyboard().Type(" ", playwright.KeyboardTypeOptions{Delay: playwright.Float(300)}) },
		func() { page.Mouse().Click(float64(rand.Intn(1200)), float64(rand.Intn(800))) },
		func() { page.Keyboard().Down("Control") },
		func() { page.Keyboard().Up("Control") },
	}

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(rand.Intn(10)+5) * time.Second)
		actions[rand.Intn(len(actions))]()
	}
}

func isBlocked(page playwright.Page) bool {
	blocked, _ := page.Locator("text='Доступ ограничен'").IsVisible()
	if blocked {
		return true
	}

	// Дополнительные проверки блокировки
	if _, err := page.Locator("h1").First().TextContent(); err != nil {
		return true
	}

	return false
}

func extractWithFallback(page playwright.Page, selectors []string, fieldName string) string {
	for _, selector := range selectors {
		if val, err := page.Locator(selector).First().TextContent(); err == nil {
			return val
		}
	}
	log.Printf("Не удалось получить %s", fieldName)
	return ""
}
