package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
)

const url = "https://www.aacfree.com"

func saveToFile(bow *browser.Browser, filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = bow.Download(w)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	bow := surf.NewBrowser()
	bow.SetUserAgent(agent.Firefox())

	err := bow.Open(url)
	if err != nil {
		saveToFile(bow, "err_log.html")
		panic(err)
	}

	// for debug
	saveToFile(bow, "aacfree.html")

	sel := bow.Find("article")
	sel.Each(func(_ int, s *goquery.Selection) {
		if id, ok := s.Attr("id"); ok {
			fmt.Println("id:", id)
		}

		if a := s.Find("h2").Find("a"); a != nil {
			if href, ok := a.Attr("href"); ok {
				fmt.Println("href:", href)
			}
		}

		if img := s.Find("img"); img != nil {
			if src, ok := img.Attr("src"); ok {
				fmt.Println("img:", src)
			}
			if alt, ok := img.Attr("alt"); ok {
				fmt.Println("name:", alt)
			}
		}

		t := s.Find("footer").Find("time")
		t.Each(func(_ int, s *goquery.Selection) {
			if class, ok := s.Attr("class"); ok {
				if strings.Contains(class, "updated") {
					if r, ok := s.Attr("datetime"); ok {
						fmt.Println("time:", r)
					}
				}
			}
		})

		fmt.Print("\n")
	})
}
