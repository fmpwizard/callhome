package main

//server
import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/read", read)
	http.HandleFunc("/send", sendMessage)
	http.HandleFunc("/ack", ackMessage)
	log.Println("Starting server ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendMessage(rw http.ResponseWriter, req *http.Request) {
	room := req.FormValue("room")
	message := req.FormValue("message")
	payload, _ := http.NewRequest("PUT", "http://127.0.0.1:4001/v2/keys/"+room, bytes.NewReader([]byte(message)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	client := &http.Client{}
	res, err := client.Do(payload)
	defer res.Body.Close()
	if err != nil {
		log.Printf("ERROR: Could not update key in etcd: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

}

//TODO: implement this one
func ackMessage(rw http.ResponseWriter, req *http.Request) {
	room := req.FormValue("room")
	message := req.FormValue("message")
	payload, _ := http.NewRequest("PUT", "http://127.0.0.1:4001/v2/keys/"+room, bytes.NewReader([]byte(message)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	client := &http.Client{}
	res, err := client.Do(payload)
	defer res.Body.Close()
	if err != nil {
		log.Printf("ERROR: Could not update key in etcd: %v", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

}

func read(rw http.ResponseWriter, req *http.Request) {
	room := req.FormValue("room")
	ret := fetchRoomMessage(room)
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(ret)
}

func fetchRoomMessage(room string) []byte {
	ret, err := http.Get("http://127.0.0.1:4001/v2/keys/" + room)
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
