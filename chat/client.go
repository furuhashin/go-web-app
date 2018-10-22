package main

import (
	"github.com/gorilla/websocket"
	"time"
)

type client struct {
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		//socket.send(msgBox.val())でテキストが送られ,ここで受ける。送受信の起点となる箇所
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string) //型アサーション
			msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		//ここでクライアントに向けてメッセージが送信される
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
