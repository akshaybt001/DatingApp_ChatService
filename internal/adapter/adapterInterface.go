package adapter

import "github.com/akshaybt001/DatingApp_ChatService/entities"

type ChatAdapterInterface interface {
	InsertMessage(msg entities.InsertIntoRoomMessage) error
	LoadMessages(roomId string) ([]entities.Message, error)
}