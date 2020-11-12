package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jeromeshapiro/geoToH3/polygon"
	geojson "github.com/paulmach/go.geojson"
	"io"
	"os"
	"strconv"
)

func main() {
	inputFile := flag.String("in", "", "input H3 file")
	outputFile := flag.String("out", "", "output GeoJSON file")
	flag.Parse()

	reader := parseInput(*inputFile)
	writer := parseWriter(*outputFile)

	exec(reader, writer)
}

func exec(reader io.Reader, writer io.Writer) {
	featureCollection := geojson.NewFeatureCollection()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		hexagon, err := strconv.ParseUint(scanner.Text(), 16, 64)
		if err != nil {
			panic(fmt.Errorf("error while parsing H3 index: %w", err))
		}

		feature := polygon.HexagonToFeature(hexagon)
		featureCollection.AddFeature(&feature)
	}

	marshalledFeatureCollection, err := featureCollection.MarshalJSON()
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprint(writer, string(marshalledFeatureCollection))
	if err != nil {
		panic(err)
	}
}

func parseInput(inputFile string) io.Reader {
	if inputFile == "" {
		return bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(inputFile)
		if err != nil {
			panic(fmt.Errorf("error while reading file: %w", err))
		}

		return file
	}
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
