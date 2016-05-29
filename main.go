package main

import (
	"fmt"
	"log"
)

func main() {
	if err := run(NewJourneyMapper()); err != nil {
		log.Fatal(err.Error())
	}
}

// Runner is the main abstraction of jobs which implement the Run method.
// Basically everything that is done by this tool eventually implements
// this interface.
type Runner interface {
	Run() error
	Name() string
}

func run(runners ...Runner) error {
	for _, r := range runners {
		log.Printf("Executing %s", r.Name())
		err := r.Run()
		if err != nil {
			return fmt.Errorf("Failed to execute runner %s: %s", r.Name(), err.Error())
		}
	}
	return nil
}
