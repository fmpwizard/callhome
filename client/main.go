package main

//client
import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	for {
		msg := string(fetchRoomMessage("Stella"))
		log.Printf("Message: %s", msg)
		err := exec.Command("espeak", "-ven+f3", "-k5", "-s120", msg).Run()
		if err != nil {
			log.Printf("ERROR: Could not read the message %s", err.Error())
		}
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
