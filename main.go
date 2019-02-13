package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"text/template"
)

// Lebski - TODO: Convert all this maps in a Database-scheme to make them persistent.
var rooms map[int]*room
var conversations map[int]*Conversation
var users map[int]*ChatUser
var usersMail map[string]*ChatUser
var templates = template.Must(template.ParseGlob("templates/*"))

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	rooms = make(map[int]*room)
	conversations = make(map[int]*Conversation)
	users = make(map[int]*ChatUser)
	usersMail = make(map[string]*ChatUser)

	user1 := newUser("test1@test.de")
	user2 := newUser("test2@test.de")
	user3 := newUser("test3@test.de")
	_, err := newConversation(user1, user2.Id)
	if err != nil {
		log.Fatal("Could not create initial conversation")
	}
	_, err = newConversation(user2, user3.Id)
	if err != nil {
		log.Fatal("Could not create initial conversation")
	}

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	router.HandleFunc("/", getLoginHandler).Methods("GET")
	router.HandleFunc("/users/{uid}/chat/{conversation}", chatHandler)
	router.HandleFunc("/users/{email}/select", selectChathandler)
	router.HandleFunc("/users/{uid}/conversation/{oppid}", newConversationHandler)
	router.HandleFunc("/users/{uid}/room/{conversation}", roomHandler)
	http.Handle("/", router)

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
