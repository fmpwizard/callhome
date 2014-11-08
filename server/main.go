package main

//server
import (
	"errors"
	"flag"
	"fmt"
	"github.com/iron-io/iron_go/mq"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//define port to listen to

var rooms = []string{"Stella", "Diego", "Hayley", "Gabriel"}

type templateData struct {
	Rooms []string
	Msg   string
}

func main() {
	var port string
	fs := http.FileServer(http.Dir("static"))
	flag.StringVar(&port, "port", "7070", "port number to bind to")
	flag.Parse()
	http.HandleFunc("/", index)
	http.HandleFunc("/read", read)
	http.HandleFunc("/crash", crash)

	http.Handle("/js/", http.StripPrefix("/js/", fs))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	log.Println("Starting server ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func index(rw http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		sendMessage(rw, req)
	} else {
		serveIndex(rw, req)
	}
}

func serveIndex(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(rw, "Could not parse html template, ERROR: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(rw, "index.html", templateData{Rooms: rooms, Msg: ""})
	if err != nil {
		log.Printf("ERROR: template gave %s", err.Error())
	}
}

func sendMessage(rw http.ResponseWriter, req *http.Request) {
	msg := req.FormValue("message")
	_ = req.ParseForm()
	for _, r := range req.Form["room"] {
		err := sendToQueue(r, msg)
		if err != nil {
			http.Error(rw, "Missing message text", http.StatusInternalServerError)
			return
		}
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(rw, "Could not parse html template, ERROR: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(rw, "index.html", templateData{Rooms: rooms, Msg: "Message was sent."})

}

func sendToQueue(room, msg string) error {
	if strings.TrimSpace(room) == "" {
		return errors.New("Missing room name")
	}
	queue := mq.New(room)
	if strings.TrimSpace(msg) == "" {
		return errors.New("Missing message text")
	}
	_, err := queue.PushString(msg)
	if err != nil {
		log.Printf("ERROR: Could not send message to queue: %v", err)
		return errors.New("Could not send message to queue")
	}
	return nil

}

/*func pushMessageToAll(msg string) {
	queue := mq.New(queueName)
}
*/
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

func crash(rw http.ResponseWriter, req *http.Request) {
	log.Println("Simulating a crash, goodbye!")
	os.Exit(1)
}
