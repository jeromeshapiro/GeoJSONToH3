package main

import (
	"fmt"
	geojson "github.com/paulmach/go.geojson"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	mergedFeatureCollection := geojson.NewFeatureCollection()
	filePaths := os.Args[1:]

	for _, filePath := range filePaths {
		input, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal("error while reading input file: %w", err)
		}

		featureCollection, err := geojson.UnmarshalFeatureCollection(input)
		if err != nil {
			log.Fatal("error while unmarshalling GeoJSON: %w", err)
		}

		for _, feature := range featureCollection.Features {
			mergedFeatureCollection.AddFeature(feature)
		}
	}

	marshalledFeatureCollection, err := mergedFeatureCollection.MarshalJSON()
	if err != nil {
		log.Fatal("error while marshalling feature collection: %w", err)
	}

	_, err = fmt.Fprint(os.Stdout, string(marshalledFeatureCollection))
	if err != nil {
		log.Fatal("error while writing to writer: %w", err)
	}
}
