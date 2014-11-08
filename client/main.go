package main

//client
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

//define server url  and port
// define room to listen to

func main() {

	var host string
	var port string
	var room string
	flag.StringVar(&port, "port", "7070", "Server port number to connect to")
	flag.StringVar(&host, "host", "127.0.0.1", "Server host (ip or domain) to connect to")
	flag.StringVar(&room, "room", "All", "The room name you are in")
	flag.Parse()

	for {
		msg := string(fetchRoomMessage(room))
		if len(strings.TrimSpace(msg)) > 0 {
			log.Printf("Message: %s", msg)
			err := exec.Command("espeak", "-ven+f3", "-k5", "-s120", msg).Run()
			if err != nil {
				log.Printf("ERROR: Could not read the message %s", err.Error())
			}
		}
		time.Sleep(2000 * time.Millisecond)
	}
}

func fetchRoomMessage(room string) []byte {
	serverUrl := fmt.Sprintf("http://%s:%s", host, port)
	ret, err := http.Get(serverUrl + "/read?room=" + room)
	if err != nil {
		log.Printf("ERROR: Could not get message for room: %s, we got error: %s", room, err.Error())
		return nil
	}
	defer ret.Body.Close()
	body, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		log.Printf("ERROR: Could not read the body response, we got error: %s", err.Error())
		return nil
	}
	return body
}
