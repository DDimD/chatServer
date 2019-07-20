package chat

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//Server manages connected clients and message forwarding
type Server struct {
	pattern      string
	broadcast    chan *OutputMessage
	addClient    chan *Client
	removeClient chan *Client
	clients      map[string]*Client
	rwMutex      sync.RWMutex
}

//NewServer create server object
func NewServer(pattern string) *Server {
	return &Server{
		pattern,
		make(chan *OutputMessage),
		make(chan *Client),
		make(chan *Client),
		make(map[string]*Client),
		sync.RWMutex{},
	}
}

//ConnectHandler handles client connection by websocket
func (srv *Server) connectHandler(w http.ResponseWriter, r *http.Request) {
	//read username from request
	err := r.ParseForm()
	if err != nil {
		log.Printf("connectHandler ParseForm err: %v", err)
	}

	username := r.FormValue("username")

	//check valid value username parameter
	if len(username) < 1 {
		log.Println("wrong username param")
		return
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	webSocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket connection upgrade error: %v", err)
		return
	}

	//check username existence
	if srv.userExist(username) {
		webSocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(4018, "username already used"))
		webSocket.Close()
		log.Printf("Username %s already used", username)
		return
	}

	client := NewClient(username, webSocket, srv)
	srv.addClient <- client
	client.Listen()
}

/*Listen creates handlers to verify that a username is busy,
connect a client by socket, handle adding and deleting clients,
 forwarding an incoming message to other clients*/
func (srv *Server) Listen() {
	http.HandleFunc("/checkUserName", srv.checkUsernameHandler)
	http.HandleFunc(srv.pattern, srv.connectHandler)

	for {
		select {
		case client := <-srv.addClient:
			srv.clients[client.name] = client
			log.Printf("user %s connected", client.name)

		case client := <-srv.removeClient:
			delete(srv.clients, client.name)
			log.Printf("client %s removed", client.name)

		case msg := <-srv.broadcast:
			srv.sendAll(msg)
		}
	}
}

func (srv *Server) sendAll(msg *OutputMessage) {
	for _, client := range srv.clients {
		if client.name != msg.UserName {
			client.send <- msg
		}
	}
}

//checkUsernameHandler ajax request to verify the use of the given username
func (srv *Server) checkUsernameHandler(rw http.ResponseWriter, req *http.Request) {
	message, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(fmt.Sprintf("username check request read error"))
	}
	username := string(message)

	if srv.userExist(username) {
		rw.WriteHeader(http.StatusTeapot)
		rw.Write([]byte("username already used"))
		return
	}

	rw.Write([]byte("username not used"))
}

// UserExist check user exist in server clients map
func (srv *Server) userExist(username string) bool {
	srv.rwMutex.RLock()
	defer srv.rwMutex.RUnlock()

	_, ok := srv.clients[username]
	return ok
}
