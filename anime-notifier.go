package main

import (
	"encoding/json"
	"fmt"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type AnimeSchedule struct {
	id    string
	time  string
	name  string
	jname string
}

type ScheduledAnimes struct {
	scheduledAnimes []AnimeSchedule
}

type Message struct {
	Name string
	Body string
	Time int64
}

type Messages struct {
	messages []Message
}

func main() {
	resp, err := http.Get("https://api-aniwatch.onrender.com/anime/schedule?date=2024-01-27")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var scheduledAnimes ScheduledAnimes
	err2 := json.Unmarshal(body, &scheduledAnimes)
	if err2 != nil {
		log.Fatalln(err)
	}

	fmt.Println(scheduledAnimes.scheduledAnimes)

	sb := string(body)
	log.Println(sb)

	http.Post("https://ntfy.sh/animenotifier", "text/plain",
		strings.NewReader(sb))

}
