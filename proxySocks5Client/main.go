package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/proxy"
)

/*
refer to http://mengqi.info/html/2015/201506062329-socks5-proxy-client-in-golang.html
*/

func Socks5Client(addr string, auth ...*proxy.Auth) (client *http.Client, err error) {

	dialer, err := proxy.SOCKS5("tcp", addr,
		nil,
		&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	)
	if err != nil {
		return
	}

	transport := &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	client = &http.Client{Transport: transport}

	return
}

func main() {
	client, err := Socks5Client("127.0.0.1:1080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resp, err := client.Get("http://www.google.com")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(b))
	}

}
