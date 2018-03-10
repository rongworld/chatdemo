package main

import (
	"github.com/gorilla/websocket"
	"fmt"
	"strconv"
	"encoding/json"
)


func (s *Client) saveOfflineMsg(m Msg) {
	strDate := fmt.Sprintf("%d", m.Date)
	data := m.Id + "-" + strDate + "|" + m.Text//redis格式： 用户名-时间戳（纳秒）|消息
	rc := RedisClient.Get()
	rc.Do("lpush", s.user.chatId, data) //键用聊天对象的用户名，聊天对象上线则查找自己的消息
	defer rc.Close()
}

func (s *Client) sendOfflineMsg(conn *websocket.Conn) {
	go s.readOfflineMsg()//读取离线消息
	for {
		select {
		case msg := <-s.user.msgChan:
			msgJson, _ := json.Marshal(msg)
			conn.WriteMessage(websocket.TextMessage, msgJson)
		}
	}
}

func (s *Client) readOfflineMsg() {
	fmt.Println("reading redis.......")
	rc := RedisClient.Get()
	defer rc.Close()
	for {
		ovalue, err := rc.Do("lpop", s.user.userId)
		value := fmt.Sprintf("%s", ovalue)
		if err != nil {
			fmt.Println(err)
			break
		}

		if ovalue == nil {
			fmt.Println("没有离线消息了")
			break
		} else {
			location1 := unicodeIndex(value,"-")

			chatId := substring(value, 0, location1)

			if chatId == s.user.chatId {
				location2 := unicodeIndex(value,"|")
				strDate := substring(value, location1+1, location2)
				date, err := strconv.ParseInt(strDate, 10, 64)
				if err != nil {
					fmt.Print(err)
					continue
				}
				text := substring(value, location2+1, len([]rune(value)))
				msg := Msg{
					Date: date,
					Text: text,
					Id:   chatId,
				}
				s.user.msgChan <- msg
			}
		}
	}
}
