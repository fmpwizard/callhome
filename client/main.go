package main

//client
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	for {
		var message RoomMessage
		json.Unmarshal(fetchRoomMessage("stella"), &message)
		log.Println(message.Node.Value)
		//log.Println(string(fetchRoomMessage("stella")))
		time.Sleep(10 * time.Millisecond)
	}
}

func fetchRoomMessage(room string) []byte {
	ret, err := http.Get("http://127.0.0.1:8080/read?room=" + room)
	if err != nil {
		log.Printf("ERROR: Could not get message for key: %s, we got error: %s", room, err.Error())
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

type RoomMessage struct {
	Action string
	Node   Etcdnode
}

type Etcdnode struct {
	Key           string
	Value         string
	ModifiedIndex int
	CreatedIndex  int
}
