package initializer

import (
	"github.com/akshaybt001/DatingApp_ChatService/internal/adapter"
	"github.com/akshaybt001/DatingApp_ChatService/internal/handler"
	"github.com/akshaybt001/DatingApp_ChatService/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func Initializer(db *mongo.Database) *handler.ChatHandlers {
	adapter := adapter.NewChatAdapter(db)
	usecase := usecase.NewChatUsecase(adapter)
	insertRoom := usecase.InsertIntoDB()
	handler := handler.NewChatHandlers(insertRoom, &usecase, "user-service:8081")
	return handler
}
