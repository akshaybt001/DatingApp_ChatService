package chat

type ChatRoom struct {
	Room map[string]*Room
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		Room: make(map[string]*Room),
	}
}