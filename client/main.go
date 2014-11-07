package main

//client
import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	for {
		log.Printf("Message: %s", string(fetchRoomMessage("Stella")))
		time.Sleep(2000 * time.Millisecond)
	}
}

func fetchRoomMessage(room string) []byte {
	ret, err := http.Get("http://127.0.0.1:7070/read?room=" + room)
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
