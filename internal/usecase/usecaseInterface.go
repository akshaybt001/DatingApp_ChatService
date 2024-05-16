package usecase

import (
	"github.com/akshaybt001/DatingApp_ChatService/entities"
	"github.com/akshaybt001/DatingApp_ChatService/internal/usecase/chat"
)

type ChatUsecaseInterface interface {
	CreateRoomifnotalreadyExists(string, chan<- entities.InsertIntoRoomMessage) (*chat.Room, []entities.Message)
}
