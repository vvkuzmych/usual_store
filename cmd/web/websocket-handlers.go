package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketConnection struct {
	*websocket.Conn
}

type WsPayload struct {
	Action      string              `json:"action"`
	Message     string              `json:"message"`
	Username    string              `json:"username"`
	MessageType string              `json:"message_type"`
	UserID      int                 `json:"user_id"`
	Conn        WebSocketConnection `json:"-"`
}

type WsJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[WebSocketConnection]string)

var wsChannel = make(chan WsPayload)

func (app *application) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	if _, exists := clients[conn]; exists {
		app.infoLog.Println("Closing duplicate WebSocket connection")
		_ = ws.Close()
		return
	}

	clients[conn] = ""
	app.infoLog.Printf("New connection from %s", r.RemoteAddr)

	go app.ListenForWS(&conn)
}

func (app *application) ListenForWS(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Println("Recovered from panic:", r)
		}
		_ = conn.Close()       // Ensure the connection is closed
		delete(clients, *conn) // Remove the client from the clients map
	}()

	var payload WsPayload
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			// Handle the "close" error more gracefully
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				app.infoLog.Println("WebSocket connection closed:", err)
			} else {
				app.errorLog.Println("Error reading from WebSocket:", err)
			}
			break
		}
		payload.Conn = *conn
		wsChannel <- payload
	}
}

func (app *application) ListenToWsChannel() {
	var response WsJsonResponse
	for {
		e := <-wsChannel
		switch e.Action {
		case "deleteUser":
			response.Action = "logout"
			response.Message = "Your account has been deleted"
			response.UserID = e.UserID
			app.broadcastToAll(response)
		default:

		}
	}
}

func (app *application) broadcastToAll(response WsJsonResponse) {
	for client := range clients {
		err := client.Conn.WriteJSON(response)
		if err != nil {
			app.errorLog.Printf("Websocker err on %s: %s", response.Action, err.Error())
			_ = client.Close()
			delete(clients, client)
		}
	}
}
