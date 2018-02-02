package main

import (
    "net/http"
    "log"
    "io"
    "fmt"
    "github.com/clevergo/captcha"
)

func main ()  {
    http.HandleFunc("/111", test)
    http.Handle("/img", captcha.Server(captcha.StdWidth, captcha.StdHeight))
    if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatal(err)
    }
}

func test(w http.ResponseWriter, r *http.Request)  {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Add("Access-Control-Allow-Headers","Content-Type")
    w.Header().Set("content-type","application/json")
    fmt.Println(r)
    io.WriteString(w, "这是从后台发送的数据")
}