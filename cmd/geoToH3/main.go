package main

import (
	"flag"
	"fmt"
	"github.com/jeromeshapiro/geoToH3/polyfill"
	geojson "github.com/paulmach/go.geojson"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

func main() {
	inputFile := flag.String("in", "", "input GeoJSON file")
	outputFile := flag.String("out", "", "output H3 file")
	resolution := flag.Int("resolution", 7, "the resolution of the output hexagons")
	strategy := flag.String("strategy", "centred", "GeoJSON to H3 polyfill algo <centred | intersecting>")
	flag.Parse()

	featureCollection := parseInput(*inputFile)
	writer := parseWriter(*outputFile)

	var polyfillStrategy polyfill.Strategy

	switch *strategy {
	case "centered":
		polyfillStrategy = polyfill.Centered
	case "intersected":
		polyfillStrategy = polyfill.Intersected
	default:
		log.Fatalf("invalid polyfill strategy: %s", *strategy)
	}

	h3Resolution := *resolution

	if h3Resolution < 0 || h3Resolution > 15 {
		log.Fatalf("invalid h3 resolution: %d", h3Resolution)
	}

	hexagons := polyfill.New(polyfillStrategy, h3Resolution).FeatureCollectionToH3(featureCollection)

	for _, hexagon := range hexagons {
		_, err := fmt.Fprintln(writer, strconv.FormatUint(hexagon, 16))
		if err != nil {
			log.Fatal("unexpected error while writing output: %w", err)
		}
	}
}

func parseInput(inputFile string) geojson.FeatureCollection {
	var input []byte
	var err error
	if inputFile == "" {
		input, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(fmt.Errorf("error while reading from STDIN: %w", err))
		}
	} else {
		input, err = ioutil.ReadFile(inputFile)
		if err != nil {
			panic(fmt.Errorf("error while reading file: %w", err))
		}
	}

	featureCollection, err := geojson.UnmarshalFeatureCollection(input)
	if err != nil {
		panic(fmt.Errorf("invalid GeoJSON input: %w", err))
	}

	return *featureCollection
}

func parseWriter(outputFile string) io.Writer {
	if outputFile == "" {
		return os.Stdout
	}

	file, err := os.Create(outputFile)
	if err != nil {
		panic(fmt.Errorf("error while writing output file: %w", err))
	}

	return file
}
