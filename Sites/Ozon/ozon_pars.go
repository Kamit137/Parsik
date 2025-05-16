package Ozon

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Price struct{
	Price string
	Date string
}
type Product struct {
	Name        string
	Price       []Price
	Img string
}


func Ozon(url string)Product {
	pythonScript := "/home/kamit/kod/Parsik/Sites/Ozon/ozon_parser.py" // Замените на фактический путь

	cmd := exec.Command("python", pythonScript, url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing Python script:", err)
		fmt.Println("Output:", string(out))
		return Product{}
	}

	// Преобразуем вывод Python-скрипта в строку JSON
	jsonString := string(out)

	// Создаем экземпляр структуры Product
	var product Product

	err = json.Unmarshal([]byte(jsonString), &product)
	if err != nil {
		fmt.Println("JSON decode error:", err)
		fmt.Println("JSON string:", jsonString) // Выводим JSON для отладки
		return Product{}
	}


	return product
}