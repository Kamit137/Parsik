package wb_pars

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
	"strconv"
	"strings"
	 Ozon "Parsik/Sites/Ozon"
)

// ------------------------------------------------------------------
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


//-----------------------------------------------------

func fail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func Wb(productURL string) Ozon.Product{
	url, err := extractCardURL(productURL, "card")
	result, err := http.Get(url)
	fail(err, "request EB fail")

	if result.StatusCode > 299 {
		fmt.Println("error ", result.StatusCode)
		os.Exit(0)
	}

	body, err := io.ReadAll(result.Body)
	fail(err, "read body")

	var jsonData JSON

	err = json.Unmarshal(body, &jsonData)
	fail(err, "json parser")

	priceFloat := float64(jsonData.Data.Products[0].SalePriceU) / 100
	priceString := fmt.Sprintf("%.2f", priceFloat)

	product := Ozon.Product{
		Name:     jsonData.Data.Products[0].Name,
		Price:    []Ozon.Price{{Price: priceString, Date: time.Now().Format("02.01.2006"),},},
		Img: ScrapeImageURLs(productURL),
	}

	return product



}

func extractCardURL(productURL string, variant string) (string, error) {
	parsedURL, err := url.Parse(productURL)
	fail(err, "некорректный URL")
	queryValues := parsedURL.Query()
	nm := queryValues.Get("nm")

	if nm == "" {
		pathParts := strings.Split(parsedURL.Path, "/")
		if len(pathParts) > 2 {
			lastPart := pathParts[len(pathParts)-2]
			nm = lastPart

		} else {
			return "", fmt.Errorf("не удалось извлечь ID товара (nm)")
		}
	}
	cardURL := fmt.Sprintf("https://card.wb.ru/cards/detail?appType=1&curr=rub&dest=-1257786&spp=27&nm=%s", nm)
	if variant == "card" {
		return cardURL, nil
	} else {
		return nm, nil
	}
}

func ScrapeImageURLs(productURL string) string {
	art, err := extractCardURL(productURL, "")
	fail(err, "fail photo")
	artI, err := strconv.Atoi(art)
	fail(err, "fail int")
	short_art := artI / 100000
	short_art2 := artI / 1000

	string_short_art := strconv.Itoa(short_art)

	string_short_art2 := strconv.Itoa(short_art2)
	basket := ""
	switch {
	case 0 <= short_art && short_art <= 143:
		basket = "01"
	case 144 <= short_art && short_art <= 287:
		basket = "02"
	case 288 <= short_art && short_art <= 431:
		basket = "03"
	case 432 <= short_art && short_art <= 719:
		basket = "04"
	case 720 <= short_art && short_art <= 1007:
		basket = "05"
	case 1008 <= short_art && short_art <= 1061:
		basket = "06"
	case 1062 <= short_art && short_art <= 1115:
		basket = "07"
	case 1116 <= short_art && short_art <= 1169:
		basket = "08"
	case 1170 <= short_art && short_art <= 1313:
		basket = "9"
	case 1314 <= short_art && short_art <= 1601:
		basket = "10"
	case 1602 <= short_art && short_art <= 1655:
		basket = "11"
	case 1656 <= short_art && short_art <= 1919:
		basket = "12"
	case 1920 <= short_art && short_art <= 2045:
		basket = "13"
	case 2046 <= short_art && short_art <= 2189:
		basket = "14"
	case 2190 <= short_art && short_art <= 2405:
		basket = "15"
	default:
		basket = "16"
	}
	img := fmt.Sprintf("https://basket-%s.wb.ru/vol%s/part%s/%s/images/big/1.webp", basket, string_short_art, string_short_art2, art)

	return img
}
