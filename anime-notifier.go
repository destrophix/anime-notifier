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
	favouriteAnimeNames := readFile()
	resp, err := http.Get("https://api-aniwatch.onrender.com/anime/schedule?date=2024-01-27")
	checkError(err)

	responseBody, err := ioutil.ReadAll(resp.Body)
	checkError(err)

	var scheduledAnimes ScheduledAnimes
	err2 := json.Unmarshal(responseBody, &scheduledAnimes)
	checkError(err2)

	favouriteAnimes := []AnimeSchedule{}

	for _, scheduledAnime := range scheduledAnimes.ScheduledAnimes {
		if slices.Contains(favouriteAnimeNames, scheduledAnime.Name) {
			favouriteAnimes = append(favouriteAnimes, scheduledAnime)
		}
	}

	data, err3 := json.Marshal(favouriteAnimes)
	checkError(err3)

	favouriteAnimesEncoded := string(data)
	log.Println(favouriteAnimesEncoded)

	http.Post("https://ntfy.sh/animenotifier", "text/plain",
		strings.NewReader(favouriteAnimesEncoded))
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
	c.AddFunc("* 17 * * *", fun)
	c.Start()
	select {}
}
