package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	URL_PREFIX = "https://trade.fenqile.com/order/detail/"
	URL_SUFFIX = ".html"
)

var inputFile = "input.txt"
var urlSlice = []string{}

func main() {
	flag.StringVar(&inputFile, "f", "input.txt", "input file")
	flag.Parse()
	fmt.Println("input file:", inputFile)

	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(line, "订单号") {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				break
			}
			if !strings.Contains(line, "<span>O") {
				continue
			}
			//fmt.Print(line)
			orderNum := strings.Split(line, "<span>")[1]
			orderNum = strings.Split(orderNum, "<")[0]
			//fmt.Println(orderNum)
			urlSlice = append(urlSlice, URL_PREFIX+orderNum+URL_SUFFIX)
		}
	}

	for _, l := range urlSlice {
		fmt.Println(l)
	}
}
