package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"http://github.com/gocolly/colly/v2"
	"github.com/PuerkitoBio/goquery"
)
func fail(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}
func main(){
cookies,err := os.ReadFile(".cookies.txt")
fail(err, "trable")
c := colly.Newcollector()
}