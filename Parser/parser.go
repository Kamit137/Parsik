package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	//"golang.org/x/net/html"
)

func Pars() {
	for _, url := range os.Args[1:] {
		links, err := findlinks(url)
		if err != nil {
			fmt.Println(os.Stderr, "Ошибка", err)
		}
		for _, link := range links {
			fmt.Println(link)
		}
	}
}

func findlinks(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as html :%v ", url, err)
	}
	return visit(nil, doc), nil
}

func visit(links []string, n *html.Node) []string {
}
