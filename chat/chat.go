package chat

import (
"fmt"
"golang.org/x/net/websocket"
"encoding/json"
)
var Users =  make(map[*websocket.Conn]string)
func WebSocket(ws *websocket.Conn) {
    var data string
    var datas Data
    for {
        //获取解析数据
        err := websocket.Message.Receive(ws, &data)
        if err != nil {
            delete(Users, ws)
            break
        }
        json.Unmarshal([]byte(data), &datas)
        if err != nil {
            fmt.Println(err)
            break
        }
        fmt.Println(datas)
        
        switch datas.Id {
        case "name":
            Users[ws] = datas.Msg
            res := `{"id":1,"msg":"系统消息：`+datas.Msg+`加入群聊"}`
            for key, _ := range Users {
                websocket.Message.Send(key, res)
            }
        case "send":
            if datas.To == "" {
                res := `{"id":3,"msg":"`+Users[ws]+":"+datas.Msg+`"}`
                for key, _ := range Users {
                    websocket.Message.Send(key, res)
                }
            } else {
                for key, val := range Users {
                    if val == datas.To {
                        res1 := `{"id":2,"msg":"你对`+val+"说:"+datas.Msg+`"}`
                        res2 := `{"id":2,"msg":"`+Users[ws]+"对你说:"+datas.Msg+`"}`
                        websocket.Message.Send(ws, res1)
                        websocket.Message.Send(key, res2)
                    }
                }
            }
        
        default:
            websocket.Message.Send(ws, "类型错误")
        }
    }
}

type Data struct {
    Id string
    Msg string
    To string
}