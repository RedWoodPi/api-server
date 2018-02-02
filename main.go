package main

import (
    "net/http"
    "log"
    "io"
    "fmt"
    "github.com/clevergo/captcha"
    "bytes"
)

type Response struct {
    Code int `json:"code"`
    Message string `json:"message"`
} 

func main ()  {
    http.HandleFunc("/111", test)
    http.Handle("/img/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
    if err := http.ListenAndServe(":80", nil); err != nil {
        log.Fatal(err)
    }
}
//处理请求发送验证码图片链接
func test(w http.ResponseWriter, r *http.Request)  {
    
    fmt.Println(r.PostForm)
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
    io.WriteString(w, buf.String())
}