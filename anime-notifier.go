package main

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type AnimeSchedule struct {
	// Id    string `json:"id"`
	Time string `json:"time"`
	Name string `json:"name"`
	// Jname string `json:"jname"`
}

type ScheduledAnimes struct {
	ScheduledAnimes []AnimeSchedule `json:"scheduledAnimes"`
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

	data, err3 := json.Marshal(scheduledAnimes)
	if err3 != nil {
		log.Fatalln(err)
	}

	sb := string(data)
	log.Println(sb)

	http.Post("https://ntfy.sh/animenotifier", "text/plain",
		strings.NewReader(sb))

}
