package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/hugo/hugolib"
)

type journey struct {
	Name string
	Slug string
}

type journeyMappingEntry struct {
	Title string `json:"title"`
}

func newJourney(path string) (*journey, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	page, err := hugolib.NewPageFrom(fp, path)
	if err != nil {
		return nil, err
	}
	return &journey{
		Slug: page.Slug,
		Name: page.Title,
	}, nil
}

func (j *JourneyMapper) findJourneys() ([]*journey, error) {
	candidates, _ := filepath.Glob("content/journey/*.md")
	result := make([]*journey, 0, len(candidates))
	for _, path := range candidates {
		jo, err := newJourney(path)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse %s: %s", path, err.Error())
		}
		j.infoLog.Printf("Found %s", path)
		result = append(result, jo)
	}
	return result, nil
}

func (j *JourneyMapper) buildJourneyMapping() error {
	target := filepath.Join("data", "journeys.json")
	journeys, err := j.findJourneys()
	if err != nil {
		return err
	}
	mapping := make(map[string]journeyMappingEntry)
	for _, j := range journeys {
		mapping[j.Slug] = journeyMappingEntry{Title: j.Name}
	}
	fp, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer fp.Close()
	err = json.NewEncoder(fp).Encode(mapping)
	if err != nil {
		return err
	}
	j.infoLog.Printf("%s updated", target)
	return nil
}

// JourneyMapper generates the data/journeys.json file out of all the journies
// stored under content/journey. This is used within the side menu and also
// facilitates linking between posts and journeys.
type JourneyMapper struct {
	infoLog *log.Logger
}

// Name returns the name of the job.
func (j *JourneyMapper) Name() string {
	return "journey mapper"
}

// Run generates the data/journeys.json file out of the available content.
func (j *JourneyMapper) Run() error {
	return j.buildJourneyMapping()
}

// NewJourneyMapper creates a new JourneyMapper instance with the default
// loggers enabled.
func NewJourneyMapper() *JourneyMapper {
	j := JourneyMapper{}
	j.infoLog = log.New(os.Stderr, "[journeymapper] ", 0)
	return &j
}
