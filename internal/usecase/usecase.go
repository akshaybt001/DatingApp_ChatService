package usecase

import (
	"fmt"
	"log"
	"strings"

	"github.com/akshaybt001/DatingApp_ChatService/entities"
	"github.com/akshaybt001/DatingApp_ChatService/internal/adapter"
	"github.com/akshaybt001/DatingApp_ChatService/internal/usecase/chat"
)

type ChatUsecase struct {
	adapter adapter.ChatAdapterInterface
	Chat    *chat.ChatRoom
}

func NewChatUsecase(adapter adapter.ChatAdapterInterface) ChatUsecase {
	return ChatUsecase{
		adapter: adapter,
		Chat:    chat.NewChatRoom(),
	}
}

func (c *ChatUsecase) CreateRoomifnotalreadyExists(roomid string, insertChan chan<- entities.InsertIntoRoomMessage) (*chat.Room, []entities.Message) {
	res, err := c.adapter.LoadMessages(roomid)
	ids := strings.Split(roomid, "")
	if err!=nil{
		log.Println("error while loading messages",err)
		res,err=c.adapter.LoadMessages(ids[1]+""+ids[0])
		if err!=nil{
			log.Println("eroor retrieving messages",err)
		}
	}
	if c.Chat.Room[roomid]==nil{
		if c.Chat.Room[ids[1]+""+ids[0]]==nil{
			fmt.Println("no message found for id ",ids[1],ids[0])
			room:=chat.NewRoom(ids[1]+""+ids[0])
			go room.Serve(insertChan)
			c.Chat.Room[ids[1]+""+ids[0]]=room
			return room,res
		}
		return c.Chat.Room[ids[1]+""+ids[0]],res
	}
	fmt.Println("room id is ",roomid)
	return c.Chat.Room[roomid],res
}
