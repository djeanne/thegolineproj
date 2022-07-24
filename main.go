package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

// First we create the structs to unmarshal our timeline to

type Timeline struct {
	XMLName    xml.Name   `xml:"timeline"`
	Version    string     `xml:"version"`
	Timetype   string     `xml:"timetype"`
	Eras       Eras       `xml:"eras"`
	Categories Categories `xml:"categories"`
	Events     Events     `xml:"events"`
	View       View       `xml:"view"`
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
	Color     string   `xml:"color"`
	EndsToday string   `xml:"ends_today"`
}

type Categories struct {
	XMLName    xml.Name   `xml:"categories"`
	Categories []Category `xml:"category"`
}

type Category struct {
	XMLName       xml.Name `xml:"category"`
	Name          string   `xml:"name"`
	Color         string   `xml:"color"`
	ProgressColor string   `xml:"progress_color"`
	DoneColor     string   `xml:"done_color"`
	FontColor     string   `xml:"font_color"`
}

type Events struct {
	XMLName xml.Name `xml:"events" json:"events"`
	Events  []Event  `xml:"event" json:"event"`
}

type Event struct {
	XMLName      xml.Name `xml:"event" json:"event"`
	Start        string   `xml:"start" json:"start"`
	End          string   `xml:"end" json:"end"`
	Text         string   `xml:"text" json:"text"`
	Progress     string   `xml:"progress" json:"progress"`
	Fuzzy        string   `xml:"fuzzy" json:"fuzzy"`
	Locked       string   `xml:"locked" json:"locked"`
	EndsToday    string   `xml:"ends_today" json:"ends_today"`
	Category     string   `xml:"category,omitempty" json:"category,omitempty"`
	Description  string   `xml:"description,omitempty" json:"description,omitempty"`
	DefaultColor string   `xml:"default_color" json:"default_color"`
	Milestone    string   `xml:"milestone" json:"milestone"`
}

type View struct {
	XMLName         xml.Name        `xml:"view"`
	DisplayedPeriod DisplayedPeriod `xml:"displayed_period"`
	//HiddenCategories HiddenCategories `xml:hidden_categories`
}

type DisplayedPeriod struct {
	XMLName xml.Name `xml:"displayed_period"`
	Start   string   `xml:"start"`
	End     string   `xml:"end"`
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

	var thisMonth int
	var thisDay int

	timeline := parseTimeline("testdata/testtimeline.timeline")
	now := time.Now()
	nowDay, nowMonth := now.Day(), int(now.Month())

	/* Pass the date that is to be checked for anniversaries
	current day is the default option
	*/

	flag.IntVar(&thisMonth, "month", nowMonth, "the month wherein the event took place")
	flag.IntVar(&thisDay, "day", nowDay, "the day whereupon the event took place")
	generateJson:= flag.Bool("json", false, "optional json output")
	addNewEventFromJson:= flag.Bool("update", false, "optional - add new event and update timeline")

	flag.Parse()

	onThisDay(timeline, thisMonth, thisDay)

	if *generateJson == true {
		makeJSON(timeline)
	}

	if *addNewEventFromJson == true {
		addEventsFromJson("testdata/testevent.json", timeline)
	}
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

func onThisDay(timeline Timeline, month int, day int) {
	now := time.Now()
	sort.Slice(timeline.Events.Events, func(i, j int) bool {
		return timeline.Events.Events[i].Start < timeline.Events.Events[j].Start
	})

	for _, event := range timeline.Events.Events {

		parsedStarts, err := time.Parse(pseudoISO, event.Start)
		if err != nil {
			log.Printf("Incorrect date format for '%s'\n", event.Text)
		}

		ago := now.Sub(parsedStarts)
		yearsAgo := int(ago.Hours() / 8760.0)

		if day == parsedStarts.Day() && time.Month(month) == parsedStarts.Month() {
			if !(strings.HasSuffix(event.Text, "'s birthday")) {
				fmt.Printf("On this day in history: %s (%s, %d years ago)\n", event.Text, parsedStarts.Format(textDate), yearsAgo)
			}
		}
	}
}

func makeJSON(timeline Timeline) {
	sort.Slice(timeline.Events.Events, func(i, j int) bool {
		return timeline.Events.Events[i].Start < timeline.Events.Events[j].Start
	})

	allEvents := timeline.Events.Events
	eventsJSON, err := json.MarshalIndent(allEvents, "", " ")
	err = ioutil.WriteFile("events.json", eventsJSON, 0644)
	if err != nil {
		log.Printf("Couldn't create json")
	}	

}

func (events *Events) AddEvent(event Event) {
    events.Events = append(events.Events, event)
}

func addEventsFromJson(eventJson string, timeline Timeline) {
	newEvent := parseNewEventFromJSON(eventJson)
	timeline.Events.AddEvent(newEvent)
	makeTimeline(timeline)
}

func parseNewEventFromJSON(eventJson string) Event {
	var event Event

	eventJsonFile, err := os.Open(eventJson)
	if err != nil {
		log.Printf("Couldn't parse json", err)
	}
	defer eventJsonFile.Close()

	jsonEvent, err := ioutil.ReadAll(eventJsonFile)
	err = json.Unmarshal(jsonEvent, &event)
	if err != nil {
		log.Printf("Couldn't parse the timeline file.")
	}

	return event
}

//TODO: validate event

func makeTimeline(timeline Timeline){
	file, err := xml.MarshalIndent(timeline, "", " ")
	err = ioutil.WriteFile("updated.timeline", file, 0644)
	if err != nil {
		log.Printf("Couldn't create timeline file")
	}
}

