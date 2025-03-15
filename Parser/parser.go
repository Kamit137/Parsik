package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var nomenclaure_id string

func initFlags() {
	flag.StringVar(&nomenclaure_id, "nm_id", "256329171", "nomenclature_id example: 256329171")
	flag.Parse()
	if nomenclaure_id == "" {
		fmt.Println("arg nm_id not found ")
		os.Exit(0)
	}
}

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

func main() {
	initFlags()

	url := "https://card.wb.ru/cards/v2/detail?appType=1&curr=rub&dest=-1255987&spp=30&ab_testing=false&lang=ru&nm=" + nomenclaure_id
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

	fmt.Println(jsonData)
}
