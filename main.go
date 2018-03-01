package main

import (
    "net/http"
    "log"
    "io"
    "fmt"
    "github.com/clevergo/captcha"
    "bytes"
    "reflect"
    "encoding/json"
    "html/template"
    "api/weather"
    "os"
    "runtime"
    "api/chat"
    "golang.org/x/net/websocket"
)

type Response struct {
    Code int `json:"code"`
    Message string `json:"message"`
}

func main ()  {
    http.HandleFunc("/weather", weatherController)
    http.HandleFunc("/", index)
    http.HandleFunc("/111", test)
    http.Handle("/websocket", websocket.Handler(chat.WebSocket))
    //验证码服务，暂时关闭
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
    http.Handle("/img/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
//主页面
func index(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles(path()+"view/index2.html")
    if err != nil {
        fmt.Println(err)
        return
    }
    
    t.Execute(w, nil)
}


//天气查询控制
func weatherController(w http.ResponseWriter, r *http.Request)  {
    err := r.ParseForm()
    fmt.Println(r)
    if err != nil {
        fmt.Println("数据解析错误")
    }
    str := weather.Weather(r.PostForm["msg"][0])
    io.WriteString(w, str)
}

//处理请求发送验证码图片链接，未启用
func test(w http.ResponseWriter, r *http.Request)  {
    err := r.ParseForm()
    if err != nil {
        fmt.Println("数据解析错误")
    }
    var js Response
    
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Add("Access-Control-Allow-Headers","Content-Type")
    w.Header().Set("content-type","application/json")
    
    d := struct {
        CaptchaId string
    }{
        captcha.New(),
    }
    var buf  bytes.Buffer
    buf.WriteString("img/")
    buf.WriteString(d.CaptchaId)
    buf.WriteString(".png")
    js.Code = 1
    js.Message = buf.String()
    fmt.Println(js)
    
    data := make(map[string]interface{})
    for i:=0; i < reflect.TypeOf(js).NumField(); i++ {
        data[reflect.TypeOf(js).Field(i).Name] = reflect.ValueOf(js).Field(i).Interface()
    }
    str, err := json.Marshal(data)
    fmt.Println(string(str))
    io.WriteString(w, string(str))
}


//根据系统判断文件路径
func path()(p string){
    osType := runtime.GOOS
    path, _ := os.Getwd()
    if osType == "windows" {
        path = path + "\\"
    }
    if osType == "linux" {
        path = path + "/"
    }
    return path
}