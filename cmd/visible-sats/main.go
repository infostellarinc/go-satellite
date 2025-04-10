package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/infostellarinc/go-satellite"
)

type State int

const (
	StateNone State = iota
	StateLine1
)

func main() {
	start := time.Now()
	inputFile := flag.String("file", "", "Input file to read (required)")
	altitude := flag.Float64("alt", 0.0, "Altitude (required)")
	longitude := flag.Float64("lon", 0.0, "Longitude (required)")
	latitude := flag.Float64("lat", 0.0, "Latitude (required)")

	// Parse flags
	flag.Parse()

	// Check if required flags are provided
	if *inputFile == "" {
		fmt.Println("All flags are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Print the coordinates
	fmt.Printf("Altitude: %.2f\n", *altitude)
	fmt.Printf("Longitude: %.2f\n", *longitude)
	fmt.Printf("Latitude: %.2f\n", *latitude)
	fmt.Printf("File: %s\n", *inputFile)

	// Open the input file
	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var coordinates = satellite.Coordinates{
		Altitude:  *altitude,
		Longitude: *longitude,
		Latitude:  *latitude,
	}

	// Create a scanner to read the file
	scanner := bufio.NewScanner(file)
	var state State
	var line1 string
	var line2 string
	tlesParsed := 0
	tleErrors := 0
	aboveHorizon := 0
	belowHorizon := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || (!strings.HasPrefix(line, "1") && !strings.HasPrefix(line, "2")) || len(line) < 60 {
			continue
		}
		switch state {
		case StateNone:
			line1 = line
			state = StateLine1
			continue
		case StateLine1:
			line2 = line
			state = StateLine1
		}
		sat, err := satellite.TLEToSat(line1, line2, satellite.GravityWGS72)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not parse TLE: %v\n", err)
			tleErrors++
			state = StateNone
			continue
		} else {
			tlesParsed++
			epoch := sat.Tle.EpochTime()
			pos, _, err := satellite.Propagate(sat, time.Now())
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not propagate satellite: %v\n", err)
				state = StateNone
				continue
			}
			lookAngles := satellite.ECIToLookAngles(pos, coordinates, satellite.JDayTime(time.Now()), sat.GravityConst)

			if lookAngles.Elevation < 0 {
				fmt.Fprintf(os.Stdout, "%v:\n\tepoch %v\n\tbelow horizon\n", sat.Tle.CatalogNumber, epoch.Format(time.RFC3339Nano))
				state = StateNone
				belowHorizon++
				continue
			}
			fmt.Fprintf(os.Stdout, "%v:\n\tepoch %v\n\tazimuth: %0.2f\n\televation: %0.2f\n", sat.Tle.CatalogNumber, epoch.Format(time.RFC3339Nano), lookAngles.Azimuth, lookAngles.Elevation)
			state = StateNone
			aboveHorizon++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
	fmt.Fprintf(os.Stdout, "Execution time: %v\nTle Parsed:Error: %v:%v\nAbove:Below horizon: %v:%v\n", time.Since(start), tlesParsed, tleErrors, aboveHorizon, belowHorizon)
}
