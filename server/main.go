package main

//server
import (
	"fmt"
	"github.com/iron-io/iron_go/mq"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/index", index)
	http.HandleFunc("/read", read)
	http.HandleFunc("/send", sendMessage)

	http.Handle("/js/", http.StripPrefix("/js/", fs))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	//http.HandleFunc("/ack", ackMessage)
	log.Println("Starting server ...")

	log.Fatal(http.ListenAndServe(":7070", nil))
}

func index(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(rw, "Could not parse html template, ERROR: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(rw, "index.html", nil)
}

var fs = http.FileServer(http.Dir("static"))

func sendMessage(rw http.ResponseWriter, req *http.Request) {
	room := req.FormValue("room")
	if strings.TrimSpace(room) == "" {
		http.Error(rw, "Missing room name", http.StatusInternalServerError)
		return
	}
	queue := mq.New(room)
	message := req.FormValue("message")
	if strings.TrimSpace(message) == "" {
		http.Error(rw, "Missing message text", http.StatusInternalServerError)
		return
	}
	_, err := queue.PushString(message)
	if err != nil {
		log.Printf("ERROR: Could not send message to queue: %v", err)
		http.Error(rw, "Could not send message to queue", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(rw, "Could not parse html template, ERROR: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(rw, "index.html", "Message was sent.")

}

func read(rw http.ResponseWriter, req *http.Request) {
	room := req.FormValue("room")
	if strings.TrimSpace(room) == "" {
		http.Error(rw, "Missing room name", http.StatusInternalServerError)
		return
	}
	queue := mq.New(room)
	info, _ := queue.Info()
	if info.Size > 0 {
		msg, err := queue.Get()
		if err != nil {
			log.Printf("ERROR: Could not read from queue: %v", err)
			http.Error(rw, "Could not read from queue", http.StatusInternalServerError)
			return
		}
		msg.Delete()
		fmt.Fprint(rw, msg.Body)
	} else {
		rw.WriteHeader(http.StatusNotModified)
	}
}
