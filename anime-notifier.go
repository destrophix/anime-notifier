package main

import (
	"bufio"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
)

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type AnimeSchedule struct {
	Time string `json:"time"`
	Name string `json:"name"`
}

type ScheduledAnimes struct {
	ScheduledAnimes []AnimeSchedule `json:"scheduledAnimes"`
}

func fun() {
	favouriteAnimes := readFile()
	resp, err := http.Get("https://api-aniwatch.onrender.com/anime/schedule?date=2024-01-27")
	checkError(err)

	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	var scheduledAnimes ScheduledAnimes
	err2 := json.Unmarshal(body, &scheduledAnimes)
	checkError(err2)

	filteredData := []AnimeSchedule{}

	for _, ele := range scheduledAnimes.ScheduledAnimes {
		if slices.Contains(favouriteAnimes, ele.Name) {
			filteredData = append(filteredData, ele)
		}
	}

	data, err3 := json.Marshal(filteredData)
	checkError(err3)

	sb := string(data)
	log.Println(sb)

	http.Post("https://ntfy.sh/animenotifier", "text/plain",
		strings.NewReader(sb))
}

func readFile() []string {
	readFile, err := os.Open("favourites.txt")
	checkError(err)

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	return fileLines
}

func main() {
	c := cron.New()
	c.AddFunc("* * * * *", fun)
	c.Start()
	select {}
}
