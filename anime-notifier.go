package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	// "fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	// "time"
	"github.com/robfig/cron/v3"
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

func fun() {
	favouriteAnimes := readFile()
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

	filteredData := []AnimeSchedule{}

	for _, ele := range scheduledAnimes.ScheduledAnimes {
		if slices.Contains(favouriteAnimes, ele.Name) {
			filteredData = append(filteredData, ele)
		}
	}

	data, err3 := json.Marshal(filteredData)
	if err3 != nil {
		log.Fatalln(err)
	}

	sb := string(data)
	log.Println(sb)

	http.Post("https://ntfy.sh/animenotifier", "text/plain",
		strings.NewReader(sb))
}

func readFile() []string {
	readFile, err := os.Open("favourites.txt")

	if err != nil {
		fmt.Println(err)
	}
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

	// fun()
	// readFile()
	select {}
}
