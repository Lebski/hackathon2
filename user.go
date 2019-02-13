package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

type ChatUser struct {
	Id               int
	Email            string
	ConversationsIds []int
}

type Conversation struct {
	Id           int
	Participants []int
	Timestamp    time.Time
	Messages     []MsgBits
	MsgCount     int
}

type MsgBits struct {
	Id        string
	SenderId  int
	Text      string
	Timestamp time.Time
}

func newConversation(user *ChatUser, userId int) (*Conversation, error) {

	secondUser, ok := getUser(userId)
	if !ok {
		log.Println("newConversation: Could not get user")
		return &Conversation{}, fmt.Errorf("newConversation: Could not get user")
	}

	conv := Conversation{
		Id:           int(uuid.New().ID()), // Lebski - TODO: Optimize datatype
		Participants: []int{user.Id, secondUser.Id},
		Timestamp:    time.Now(),
		Messages:     []MsgBits{},
	}
	//user.Conversations = append(user.Conversations, &conv)
	//secondUser.Conversations = append(secondUser.Conversations, &conv)
	user.ConversationsIds = append(user.ConversationsIds, conv.Id)
	secondUser.ConversationsIds = append(secondUser.ConversationsIds, conv.Id)
	conversations[conv.Id] = &conv
	return &conv, nil
}

func newUser(userMail string) *ChatUser {
	user := ChatUser{
		Id:               int(uuid.New().ID()), // Lebski - TODO: Optimize datatype
		Email:            userMail,
		ConversationsIds: []int{},
	}

	fmt.Printf("Creating new user: %s --> %d\n", userMail, user.Id)
	users[user.Id] = &user
	usersMail[user.Email] = &user
	return &user
}

func getUser(userId int) (*ChatUser, bool) {
	user, ok := users[userId]
	if !ok {
		fmt.Println("User not found", userId)
		return &ChatUser{}, false
	}
	fmt.Println("User found:", userId)
	return user, true
}

func (u *ChatUser) getCounterpart(conv *Conversation) (*ChatUser, error) {
	if p1, ok := getUser(conv.Participants[0]); ok {
		if u.Id != p1.Id {
			return p1, nil
		} else {
			if p2, ok := getUser(conv.Participants[1]); ok {
				return p2, nil
			}
		}
	}
	return &ChatUser{}, fmt.Errorf("getCounterpart: Could not load users")
}

func getUserByMail(userMail string) (*ChatUser, bool) {
	user, ok := usersMail[userMail]
	if !ok {
		fmt.Println("User not found", userMail)
		return &ChatUser{}, false
	}
	fmt.Println("User found:", userMail)
	return user, true
}

func getConversation(convId int) (*Conversation, bool) {
	conv, ok := conversations[convId]
	if !ok {
		fmt.Println("Conversation not found", convId)
		return &Conversation{}, false
	}
	fmt.Println("Conversation found:", convId)
	return conv, true
}

func (conv *Conversation) addMessage(senderId int, text string) *MsgBits {
	msg := MsgBits{
		Id:        uuid.New().String(),
		SenderId:  senderId,
		Text:      text,
		Timestamp: time.Now(),
	}
	conv.MsgCount++
	fmt.Println("New Message created", conv.Id)
	conv.Messages = append(conv.Messages, msg)
	fmt.Println("New Message created", msg)
	return &msg
}
