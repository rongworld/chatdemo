package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"time"
	"encoding/json"
)



type Msg struct {
	Date int64  `json:"date"`
	Text string `json:"msg"`
	Id   string `json:"chatId"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 1 * time.Second,
}

type Client struct {
	holder *Holder

	conn *websocket.Conn

	msgChan chan []byte

	user *User
}

type User struct {
	userId  string
	chatId  string
	msgChan chan Msg
}

func (s *Client) sendMessage() { //服务器向用户发送消息
	for {
		select {
		case text := <-s.msgChan:
			chatIdServer, ok := s.holder.clients[s.user.chatId]
			msg := Msg{
				Date: time.Now().UnixNano(),
				Id:   s.user.userId, //向对方发送聊天对象id
				Text: string(text[:]),
			}

			if ok {
				//在线则直接发送信息
				//chatIdServer.conn.WriteMessage(websocket.TextMessage, text)
				msgJson, err := json.Marshal(msg)
				if err != nil {
					fmt.Println(err)
					return
				} else {
					chatIdServer.conn.WriteMessage(websocket.TextMessage, msgJson)
				}

			} else {
				s.saveOfflineMsg(msg)
			}

		}
	}
}

func (s *Client) readMessage() { //服务器读取该用户发来的消息

	for {
		_, data, err := s.conn.ReadMessage()
		if err != nil {
			fmt.Println(s.user.userId, "下线")
			s.downLine() //出现错误下线
			break
		}

		s.msgChan <- data
	}
}

func (s *Client) downLine() {
	s.conn.Close()
	delete(s.holder.clients, s.user.userId)
}

func serverWs(h *Holder, w http.ResponseWriter, r *http.Request) {

	token, _ := getToken(r)
	//chatId := fmt.Sprintf("%s", token.Header["chatId"])
	//userId := fmt.Sprintf("%s", token.Header["userId"])

	userId := getIdFromClaims("id", token.Claims)

	fmt.Println("userId"+userId)

	if userId == "" {
		return
	}
	chatId := getChatID(userId)


	if chatId == "" {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil) //新建conn

	if err != nil {
		fmt.Print(err)
		return
	}

	u := User{
		userId:  userId,
		chatId:  chatId,
		msgChan: make(chan Msg, 15),
	}

	client := &Client{//新建服务端

		holder: h,

		conn: conn,

		msgChan: make(chan []byte),

		user: &u,
	}

	h.clients[userId] = client              //上线，加入在线列表
	keepAlive(client, 700*time.Millisecond) //判断客户端是否掉线
	go client.sendOfflineMsg(conn)          //检测是否有离线消息并发送离线消息
	go client.readMessage()                 //读取用户发来的消息
	go client.sendMessage()                 //向用户发送消息
	fmt.Println(userId, "上线")
}

func keepAlive(c *Client, timeout time.Duration) {
	lastResponse := time.Now()
	c.conn.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		for {
			err := c.conn.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				c.downLine()
				return
			}
		}
	}()
}
