package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	fmt.Println("Path:", r.URL.Path)
	fmt.Println("Scheme:", r.URL.Scheme)
	r.ParseForm()
	fmt.Println("Form:", r.Form)

	fmt.Fprintf(w, "this is the root path!") //这个写入到w的是输出到客户端的
}

func handleMan(w http.ResponseWriter, r *http.Request) {
	fmt.Println("---------------------------")
	fmt.Println("Path:", r.URL.Path)
	fmt.Println("Scheme:", r.URL.Scheme)
	r.ParseForm()
	fmt.Println("Form:", r.Form)

	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("request body:\n", string(body))

	fmt.Fprintf(w, "this is the Man path!")
}

func handleLogin(res http.ResponseWriter, req *http.Request) {
	fmt.Println("method:", req.Method) //获取请求的方法
	if req.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(res, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		req.ParseForm()
		fmt.Println("username:", req.PostForm["username"]) // or req.Form[]
		fmt.Println("password:", req.PostForm["password"])
	}
}

func main() {
	//设置访问的路由
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/man", handleMan)
	http.HandleFunc("/login", handleLogin)

	err := http.ListenAndServe(":8066", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
