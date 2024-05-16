package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/akshaybt001/DatingApp_ChatService/entities"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	ClientID string
	Name     string
	Room     *Room
}

func NewClient(conn *websocket.Conn, clientID, name string, room *Room) *Client {
	return &Client{
		Conn:     conn,
		ClientID: clientID,
		Name:     name,
		Room:     room,
	}
}

func (client *Client) Serve(msgs []entities.Message) {
	client.Room.JoinChan <- client
	defer func() {
		client.Room.LeaveChan <- client
		client.Conn.Close()
	}()
	for _, v := range msgs {
		jsonData, err := json.Marshal(v)
		if err != nil {
			log.Println("error at sending ", err)
			continue
		}
		if err := client.Conn.WriteMessage(v.Type, jsonData); err != nil {
			log.Println("error at sending ", err)
			continue
		}
	}
	for {
		msgtype,p,err:=client.Conn.ReadMessage()
		if err!=nil{
			log.Println("error happened , closing connection")
			break
		}
		message:=entities.Message{Type: msgtype,Message: string(p),Time: time.Now(),Name: client.Name}
		client.Room.Broadcast<-message
		log.Printf("message recieved from %s",client.ClientID)
	}
}
