package Ozon

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type ProductInfo struct {
	Name     string  `json:"name"`
	Price    string `json:"price"`
	ImageURL string  `json:"imageUrl"`
}

func Ozon(url string)ProductInfo {
 // url := "https:///www.ozon.ru/product/sredstvo-dlya-mytya-posudy-synergetic-detskih-igrushek-c-aromatom-granata-5-l-antibakterialnoe-1436053626/"
	pythonScript := "/home/kamit/kod/Parsik/Sites/Ozon/ozon_parser.py" // Замените на фактический путь

	cmd := exec.Command("python", pythonScript, url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing Python script:", err)
		fmt.Println("Output:", string(out))
		return ProductInfo{}
	}

	// Преобразуем вывод Python-скрипта в строку JSON
	jsonString := string(out)

	// Создаем экземпляр структуры ProductInfo
	var product ProductInfo

	// Распарсиваем JSON в структуру ProductInfo
	if err := json.Unmarshal([]byte(jsonString), &product); err != nil {
		fmt.Println("JSON decode error:", err)
		fmt.Println("JSON string:", jsonString) // Выводим JSON для отладки
		return ProductInfo{}
	}


	return product
}