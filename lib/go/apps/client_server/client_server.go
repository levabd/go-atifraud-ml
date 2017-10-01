package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
	"github.com/gorilla/websocket"
	"net"
	"fmt"
	"net/http"
	"strings"
	"io/ioutil"
)

var (
	// flagPort is the open port the application listens on
	results     []string
	headerQueue []string
)

func createConnection(addr1 *string) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr1, Path: "/receive-header"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:

			if len(headerQueue) > 0 {
				for i := 0; i < len(headerQueue); i++ {
					err = c.WriteMessage(websocket.TextMessage, []byte(headerQueue[i]))
					if err != nil {
						log.Println("write:", err)
						return
					}
					headerQueue = append(headerQueue[:i], headerQueue[i+1:]...)
				}
			}

		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}

func checkIp(ip string) bool {
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		fmt.Printf("%v is not an IPv4 address\n", trial)
		return false
	}

	return true
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
		}

		headerQueue=append(headerQueue, string(body))
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

func main() {
	if len(os.Args) != 3 {
		panic("Not enough arguments ro run client. " +
			"You must provide client server port as a first argument and ip:port to send as a second argument")
		log.Println("Not enough arguments ro run client. " +
			"You must provide client server port as a first argument and ip:port to send as a second argument")
		return
	}

	ip := strings.Split(os.Args[2], ":", )

	if ip[0] != "localhost" && !checkIp(ip[0]) {
		panic("Invalid IP address provided")
	}

	if ip[0] == "" || ip[1] == "" {
		panic("Empty receiving ip/port  provided")
	}

	if os.Args[1] == "" {
		panic("Empty client server port provided")
	}

	go func() {

		flagPort:= flag.String("port",  os.Args[1], "Port to listen on")
		// set post request handler
		results = append(results, time.Now().Format(time.RFC3339))
		mux := http.NewServeMux()
		mux.HandleFunc("/send-header", PostHandler)
		log.Printf("listening on port %s", *flagPort)
		log.Fatal(http.ListenAndServe(":" + *flagPort, mux))
	}()

	// set WebSocket connection main server
	addr1 := *flag.String("addr", os.Args[2], "http service address")
	createConnection(&addr1)
}
