package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"strings"
	"github.com/buger/jsonparser"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func handleHeader(response string) bool {

	start_time := time.Now()
	start_nanosecond := start_time.Nanosecond()
	println(start_time.Minute(), start_time.Second(), start_nanosecond)

	main_data, value_data, order_data := HandleLogLine(response)

	end := time.Now()
	println(end.Minute(), end.Second(), end.Nanosecond(), end.Nanosecond()-start_nanosecond)

	fmt.Printf("%s\n", main_data)
	fmt.Printf("%+v\n", value_data)
	fmt.Printf("%+v\n", order_data)

	return true
}

func HandleLogLine(line string) (string, map[string]interface{}, map[string]interface{}) {
	var (
		elements                           = strings.SplitN(line, ",", 3)
		value_row   map[string]interface{} = make(map[string]interface{})
		ordered_row map[string]interface{} = make(map[string]interface{})
	)
	if line == "" {
		return "",  value_row, ordered_row
	}

	json_to_parse := strings.Replace(string(line), " ", "", -1)
	json_to_parse = strings.TrimPrefix(strings.TrimSuffix(json_to_parse, ""), "'")

	if len(elements[1]) > 0 {
		data := []byte(json_to_parse)
		i := 0
		jsonparser.ObjectEach(data, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if string(key) == "ip"{
				return nil
			}
			ordered_row[ string(key)] = i
			value_row[ string(key) ] = string(value)
			i = i + 1
			return nil
		})

		// define crawler in User-Agent
		if user_agent, ok := value_row["User-Agent"].(string); ok {
			return  user_agent,  value_row, ordered_row
		} else {
			return   "",  value_row, ordered_row
		}
	}
	return  "",  value_row, ordered_row
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		handleHeader(string(message))

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/receive-header", handleRequest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
