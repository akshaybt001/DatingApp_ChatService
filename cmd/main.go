package main

import (
	"log"
	"os"

	"github.com/akshaybt001/DatingApp_ChatService/db"
	"github.com/akshaybt001/DatingApp_ChatService/initializer"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("error on loading env")
	}
	addr := os.Getenv("DB_KEY")
	db, err := db.InitMongoDB(addr)
	if err != nil {
		log.Fatal("error on connecting db")
	}
	handler := initializer.Initializer(db)
	handler.Start()
}
