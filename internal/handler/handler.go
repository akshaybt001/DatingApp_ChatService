package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/akshaybt001/DatingApp_ChatService/entities"
	"github.com/akshaybt001/DatingApp_ChatService/internal/usecase"
	"github.com/akshaybt001/DatingApp_ChatService/internal/usecase/chat"
	"github.com/akshaybt001/DatingApp_proto_files/pb"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

type ChatHandlers struct {
	InsertChannel chan<- entities.InsertIntoRoomMessage
	Usecase       usecase.ChatUsecaseInterface
	UserConn      pb.UserServiceClient
	Upgrader      websocket.Upgrader
}

func NewChatHandlers(insertChannel chan<- entities.InsertIntoRoomMessage, usecase usecase.ChatUsecaseInterface, userAddr string) *ChatHandlers {
	userRes, _ := grpc.Dial(userAddr, grpc.WithInsecure())
	return &ChatHandlers{
		InsertChannel: insertChannel,
		Usecase:       usecase,
		UserConn:      pb.NewUserServiceClient(userRes),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (c *ChatHandlers) Handler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Sec-WebSocket-Extensions")
	userId1 := r.Header.Get("userId1")
	userId2 := r.Header.Get("userId2")
	recieverId := r.Header.Get("recieverId")
	var roomId string
	if userId1 != "" && recieverId != "" {
		roomId = userId1 + "" + recieverId
	} else if userId2 != "" && recieverId != "" {
		roomId = recieverId + "" + userId2
	} else {
		http.Error(w, "please provide valid headers", http.StatusBadRequest)
		return
	}
	UserData1, err := c.UserConn.GetUserData(context.Background(), &pb.GetUserById{Id: userId1})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := UserData1.Name
	UserId := userId1
	if userId2 != "" {
		userData2, err := c.UserConn.GetUserData(context.Background(), &pb.GetUserById{Id: userId2})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name = userData2.Name
	}

	conn, err := c.Upgrader.Upgrade(w, r, r.Header)
	if err != nil {
		log.Println("error while upgrading ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	room, msgs := c.Usecase.CreateRoomifnotalreadyExists(roomId, c.InsertChannel)
	client := chat.NewClient(conn, UserId, name, room)
	client.Serve(msgs)

}

func (chat *ChatHandlers) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", chat.Handler)
	log.Println("listening on port 8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		fmt.Println(err.Error())
	}
}
