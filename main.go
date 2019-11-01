package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	//"strings"
	"time"
)

// First we create the structs to unmarshal our timeline to

type Timeline struct {
	XMLName    xml.Name   `xml:"timeline"`
	Version    string     `xml:version`
	Timetype   string     `xml:timetype`
	Eras       Eras       `xml:eras`
	Categories Categories `xml:categories`
	Events     Events     `xml:events`
	View       View       `xml:view`
}

type Eras struct {
	XMLName xml.Name `xml:"eras"`
	Eras    []Era    `xml:"era"`
}

type Era struct {
	XMLName   xml.Name `xml:"era"`
	Name      string   `xml:"name"`
	Start     string   `xml:"start"`
	End       string   `xml:"end"`
	Color     string   `xml:color`
	EndsToday string   `xml:ends_today`
}

type Categories struct {
	XMLName    xml.Name   `xml:"categories"`
	Categories []Category `xml:"category"`
}

type Category struct {
	XMLName       xml.Name `xml:"category"`
	Name          string   `xml:"name"`
	Color         string   `xml:color`
	ProgressColor string   `xml:progress_color`
	DoneColor     string   `xml:done_color`
	FontColor     string   `xml:font_color`
}

type Events struct {
	XMLName xml.Name `xml:"events"`
	Events  []Event  `xml:"event"`
}

type Event struct {
	XMLName      xml.Name `xml:"event"`
	Start        string   `xml:"start"`
	End          string   `xml:"end"`
	Text         string   `xml:"text"`
	Progress     string   `xml:progress`
	Fuzzy        string   `xml:fuzzy`
	Locked       string   `xml:locked`
	EndsToday    string   `xml:ends_today`
	Category     string   `xml:"category"`
	Description  string   `xml:"description"`
	DefaultColor string   `xml:default_color`
	Milestone    string   `xml:milestone`
}

type View struct {
	XMLName         xml.Name        `xml:view`
	DisplayedPeriod DisplayedPeriod `xml:displayed_period`
	//HiddenCategories HiddenCategories `xml:hidden_categories`
}

type DisplayedPeriod struct {
	XMLName xml.Name `xml:displayed_period`
	Start   string   `xml:start`
	End     string   `xml:end`
}

// We do not need the hidden categories as of now

/*type HiddenCategories struct {
	XMLName xml.Name `xml:hidden_categories`
	Names []Name `xml:name`
	}
}*/

// Need to define the time/date layout used in the project for cleaner output

const (
	pseudoISO = "2006-01-02 15:04:05"
	textDate  = "January 2, 2006"
)

func main() {

	timeline := parseTimeline("testdata/testtimeline.timeline")
	onThisDay(timeline)
}

// Parse the timeline file created by the software

func parseTimeline(timefile string) Timeline {
	var goline Timeline
	xmlTimeline, err := os.Open(timefile)
	if err != nil {
		log.Println(err)
	}
	defer xmlTimeline.Close()

	readTimeline, _ := ioutil.ReadAll(xmlTimeline)

	err = xml.Unmarshal(readTimeline, &goline)
	if err != nil {
		log.Printf("Couldn't parse the timeline file.")
	}

	return goline
}

/* Determine whether on a given day it is the anniversary of any Timeline event
and, if yes, calculate how much time has elapsed in years */

func onThisDay(timeline Timeline) {

	now := time.Now()

	sort.Slice(timeline.Events.Events, func(i, j int) bool {
		return timeline.Events.Events[i].Start < timeline.Events.Events[j].Start
	})

	for _, event := range timeline.Events.Events {

		parsedStarts, err := time.Parse(pseudoISO, event.Start)
		if err != nil {
			log.Printf("Incorrect date format for '%s'\n", event.Text)
		}
		/*parsedEnds, err := time.Parse(pseudoISO, event.End)
		if err != nil {
			log.Printf("Incorrect date format for '%s'\n", event.Text)
		}*/

		ago := now.Sub(parsedStarts)
		yearsAgo := int(ago.Hours() / 8760.0)

		if now.Day() == parsedStarts.Day() && now.Month() == parsedStarts.Month() {
			fmt.Printf("On this day in history: %s (%s, %d years ago)\n", event.Text, parsedStarts.Format(textDate), yearsAgo)
		}
	}
}
