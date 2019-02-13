# Workgroup 4 - Task 1
https://github.com/Lebski/go-chat

### go-chat

go-chat is a realtime chat example ebedded in a rest server. 
The application is extremely easy to use, so you can get started right away and you don't have to install a database or anything like that.
Features:

  - Persistent Conversations
  - Users with accounts 
  - Websockets for realtime interaction

###  :warning: Disclaimer! :warning:

##### Security
No security features are supported.
**But:** The chat application is desigend to be simplistic so you can integrate it into your existing project without changing to much. 
I recommend using 
  - Encrypted Coockies :cookie:
  - Middleware for authentication
  - In-memory sessions 

Please make sure you **do not use URL parameters as authentication** or resource allocation. 
_further links:  [Gorilla securecookie](http://www.gorillatoolkit.org/pkg/securecookie), [Gorilla securecookie](http://www.gorillatoolkit.org/pkg/sessions), [redis](https://redis.io/)_

##### Databases 
I do not use any database but store everything inside memory, which means that everything (including chat history) ist lost when you restart the application. 
It should be easy to replace the maps with the database of your choice. 

###  Usage

```sh
$ go get github.com/gorilla/mux
$ go get github.com/gorilla/websocket
$ go get github.com/google/uuid
$ go build . 
$ ./go-chat
```