package main

import (
    "net/http"
    "log"
    "io"
    "fmt"
    "github.com/clevergo/captcha"
    "bytes"
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
    d := struct {
        CaptchaId string
    }{
        captcha.New(),
        }
    var buf  bytes.Buffer
    buf.WriteString("http://101.132.118.202/img/")
    buf.WriteString(d.CaptchaId)
    buf.WriteString(".png")
    fmt.Print(buf.String())
    io.WriteString(w, buf.String())
}