package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/akshaybt001/DatingApp_ChatService/entities"
)

type Room struct {
	ID        string
	JoinChan  chan *Client
	LeaveChan chan *Client
	Broadcast chan entities.Message
	Clients    map[string]*Client
}

func NewRoom(id string) *Room {
	return &Room{
		ID:        id,
		JoinChan:  make(chan *Client),
		LeaveChan: make(chan *Client),
		Broadcast: make(chan entities.Message),
		Clients:    make(map[string]*Client),
	}
}

type Register struct {
	Message string    `json:"Message"`
	Time    time.Time `json:"Time"`
}

func (room *Room) Serve(insertChan chan<-entities.InsertIntoRoomMessage){
	defer func(){
		close(room.JoinChan)
		close(room.LeaveChan)
		close(room.Broadcast)
	}()
	for{
		select {
		case client:=<-room.JoinChan:
			for _,v:=range room.Clients{
				reg:=Register{
					Time: time.Now(),
					Message:fmt.Sprintf("%s is online",client.Name),
				}
				if err:=v.Conn.WriteJSON(reg);err!=nil{
					log.Println("error happened at sending ",err)
					continue
				}
			}
			room.Clients[client.ClientID]=client
		case client:=<-room.LeaveChan:
			for _,v:=range room.Clients{
				unReg:=Register{
					Time: time.Now(),
					Message: fmt.Sprintf("%s is offline ",client.Name),
				}
				if err:=v.Conn.WriteJSON(unReg);err!=nil{
					log.Println("error at sending ",err)
					continue
				}
		
			}
			delete(room.Clients,client.ClientID)
		case message:=<-room.Broadcast:
			for _,v:=range room.Clients{
				jsonData,err:=json.Marshal(message)
				if err!=nil{
					log.Println("error at senting ")
					continue
				}
				if err:=v.Conn.WriteMessage(message.Type,jsonData);err!=nil{
					log.Println("error at sending")
					continue
				}
			}
			msg:=entities.InsertIntoRoomMessage{
				RoomID: room.ID,
				Messages: message,
			}
			fmt.Println("message sent ",msg)
			insertChan<-msg
		}
	}
}