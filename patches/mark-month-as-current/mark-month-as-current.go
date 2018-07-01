package main

import (
	"fmt"
	"time"

	"github.com/animenotifier/arn"
	"github.com/fatih/color"
)

func main() {
	color.Yellow("Marking current month as current")

	defer color.Green("Finished")
	defer arn.Node.Close()

	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()

	for anime := range arn.StreamAnime() {
		animeStart := anime.StartDateTime()

		if anime.Status != "current" && animeStart.Year() == currentYear && animeStart.Month() == currentMonth {
			fmt.Println(anime)
			anime.Status = "current"
			anime.Save()
		}
	}

	time.Sleep(1 * time.Second)
}
