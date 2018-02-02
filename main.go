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
    d := struct {
        CaptchaId string
    }{
        captcha.New(),
        }
    io.WriteString(w, d."http://101.132.118.202/img")
}