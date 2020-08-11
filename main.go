package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	name = flag.String("name", "", "name filter")
)

type LocationHistory struct {
	TimelineObjects []struct {
		PlaceVisit Visit `json:"placeVisit"`
	} `json:"timelineObjects"`
}

type Visit struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Duration struct {
		StartTimestampMS string `json:"startTimestampMs"`
		EndTimestampMS   string `json:"endTimestampMs"`
	} `json:"duration"`
	ChildVisits     []Visit `json:"childVisits"`
	PlaceConfidence string  `json:"placeConfidence"`
}

type Filter struct {
	NameContains string
	Confidence   string
}

type Result struct {
	Name            string
	Start           time.Time
	End             time.Time
	PlaceConfidence string
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("path to file(s) not provided")
		os.Exit(1)
	}
	paths := flag.Args()

	filter := Filter{NameContains: *name}
	if filter == (Filter{}) {
		log.Fatal("no filters provided")
	}

	var allVisits []Visit
	for _, path := range paths {
		locationHistory, err := parsePath(path)
		if err != nil {
			log.Fatal(err)
		}
		visits := flattenVisits(locationHistory)
		allVisits = append(allVisits, visits...)
	}
	results, err := filterVisits(allVisits, filter)
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func parsePath(path string) (locationHistory LocationHistory, err error) {
	f, err := os.Open(path)
	if err != nil {
		return locationHistory, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return locationHistory, err
	}
	err = json.Unmarshal(b, &locationHistory)
	if err != nil {
		return locationHistory, err
	}
	return locationHistory, nil
}

func flattenVisits(history LocationHistory) (visits []Visit) {
	for _, timelineObject := range history.TimelineObjects {
		visit := timelineObject.PlaceVisit
		if visit.Location.Name == "" {
			continue
		}
		for _, childVisit := range visit.ChildVisits {
			visits = append(visits, childVisit)
		}
		visit.ChildVisits = nil
		visits = append(visits, visit)
	}
	return visits
}

func filterVisits(visits []Visit, filter Filter) (results []Result, err error) {
	for _, v := range visits {
		if filter.NameContains != "" {
			if !strings.Contains(strings.ToLower(v.Location.Name), strings.ToLower(filter.NameContains)) {
				continue
			}
		}
		if filter.Confidence != "" && v.PlaceConfidence != filter.Confidence {
			continue
		}

		result := Result{
			Name:            v.Location.Name,
			PlaceConfidence: v.PlaceConfidence,
		}
		startTimeInt, err := strconv.ParseInt(v.Duration.StartTimestampMS, 10, 64)
		if err == nil {
			result.Start = time.Unix(0, int64(time.Millisecond)*startTimeInt)
		}
		endTimeInt, err := strconv.ParseInt(v.Duration.EndTimestampMS, 10, 64)
		if err == nil {
			result.End = time.Unix(0, int64(time.Millisecond)*endTimeInt)
		}
		results = append(results, result)
	}
	return results, nil
}
