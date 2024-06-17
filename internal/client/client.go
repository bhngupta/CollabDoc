package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func StartClient() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	fmt.Printf("Connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	for i := 0; i < 10; i++ {
		err := c.WriteJSON(map[string]interface{}{"message": fmt.Sprintf("Hello %d", i)})
		if err != nil {
			log.Println("Write error:", err)
			return
		}
	}
}
