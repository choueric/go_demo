package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	URL_SUFFIX = ".html"
)

var (
	inputFile = "input.txt"
	urlSlice  = []string{}
	config    *Config
)

func main() {
	flag.StringVar(&inputFile, "f", "input.txt", "input file")
	flag.Parse()
	fmt.Println("input file:", inputFile)

	config = getConfig()

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
			urlSlice = append(urlSlice, config.UrlPrefix+orderNum+URL_SUFFIX)
		}
	}

	for _, l := range urlSlice {
		err := execCmd(l)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

// execute command directly.
func execCmd(url string) error {
	fmt.Println("open:", url)
	cmd := exec.Command(config.Chrome, url)
	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
