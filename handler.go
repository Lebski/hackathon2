package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func getLoginHandler(w http.ResponseWriter, req *http.Request) {
	data := map[string]interface{}{}
	templates.ExecuteTemplate(w, "index", data)
}

func selectChathandler(w http.ResponseWriter, r *http.Request) {
	// get with "user" (userid) and "email"
	vars := mux.Vars(r)

	user, ok := getUserByMail(vars["email"])
	if !ok {
		user = newUser(vars["email"])
	}

	type conversationSelect struct {
		ConversationId int
		UserEmail      string
	}

	var userConversations []conversationSelect

	for _, convId := range user.ConversationsIds {
		if conv, ok := getConversation(convId); ok {
			if counterpart, err := user.getCounterpart(conv); err == nil {
				userConversations = append(userConversations, conversationSelect{conv.Id, counterpart.Email})
			}
		}
	}

	data := map[string]interface{}{
		"userId":        user.Id,
		"conversations": userConversations,
		// Lebski - TODO: USERS IS ONLY FOR DEMONSTRATION PURPOSES
		"users": users,
	}
	templates.ExecuteTemplate(w, "selectChat", data)
}

func newConversationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["uid"])
	if err != nil {
		log.Fatal("User ID could not be read")
	}

	opponentId, err := strconv.Atoi(vars["oppid"])
	if err != nil {
		log.Fatal("User ID could not be read")
	}

	user, ok := getUser(userId)
	if !ok {
		log.Println("User ID was not correct")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	conv, err := newConversation(user, opponentId)
	if err != nil {
		log.Println("Could not create Conversation")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	redirectURL := "/users/" + strconv.Itoa(user.Id) + "/chat/" + strconv.Itoa(conv.Id)
	http.Redirect(w, r, redirectURL, http.StatusFound)

}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Lebski - TODO: In production store conversationId in-memory or in encrypted cookie
	convId, err := strconv.Atoi(vars["conversation"])
	if err != nil {
		log.Fatal("Conversation ID could not be read")
	}

	userId, err := strconv.Atoi(vars["uid"])
	if err != nil {
		log.Fatal("User ID could not be read")
	}

	user, ok := getUser(userId)
	if !ok {
		log.Println("User ID was not correct")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	conv, ok := getConversation(convId)
	if !ok {
		log.Println("Conversation ID was not correct")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	opponent, err := user.getCounterpart(conv)
	if err != nil {
		log.Println("Conversation ID was not correct")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := map[string]interface{}{
		"Host":         r.Host,
		"convId":       convId,
		"userId":       userId,
		"userMail":     user.Email,
		"opponentMail": opponent.Email,
	}

	if conv.MsgCount > 0 {
		var lastMessages []MsgBits
		if conv.MsgCount > 10 {
			lastMessages = conv.Messages[conv.MsgCount-10:]
		} else {
			lastMessages = conv.Messages
		}

		msgBytes, err := json.Marshal(lastMessages)
		if err != nil {
			return
		}

		data["messages"] = string(msgBytes)
	}

	templates.ExecuteTemplate(w, "chat", data)

}

func roomHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)

	// Lebski - TODO: In production store conversationId in-memory or in encrypted cookie
	convId, err := strconv.Atoi(vars["conversation"])
	if err != nil {
		log.Println("Conversation ID could not be read")
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	userId, err := strconv.Atoi(vars["uid"])
	if err != nil {
		log.Println("User ID could not be read")
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	user, ok := getUser(userId)
	if !ok {
		log.Println("User ID was not correct")
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	conv, ok := getConversation(convId)
	if !ok {
		log.Println("Conversation ID was not correct")
		http.Redirect(w, req, "/", http.StatusFound)
		return
	}

	selectedRoom, ok := roomExists(convId)
	if !ok {
		selectedRoom = newRoom(convId)
		go selectedRoom.run()
	}

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP-:", err)
		return
	}

	client := &client{
		user:         user,
		conversation: conv,
		socket:       socket,
		send:         make(chan []byte, messageBufferSize),
		room:         selectedRoom,
	}
	selectedRoom.join <- client
	defer func() { selectedRoom.leave <- client }()
	go client.write()
	client.read()
}
