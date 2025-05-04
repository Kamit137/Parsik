package Ozon

import (
	"fmt"
	"os/exec"
)
type ProductInfo struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	ImageURL string  `json: imageUrl`
}
func Ozon(url string)ProductInfo {
	cmd := exec.Command("python", "Sites/Ozon/ozon_parser.py", url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	product := ProductInfo{
		Name: string(out[0]),
		Price: string(out[1]),

		ImageURL: string(out[2]),
	}
	return product
}
