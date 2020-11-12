package polyfill

import (
	geojson "github.com/paulmach/go.geojson"
)

type Strategy int

const (
	Centered = iota
	Intersected
)

type Polyfill interface {
	FeatureCollectionToH3(featureCollection geojson.FeatureCollection) []uint64
}

func New(strategy Strategy, resolution int) Polyfill {
	switch strategy {
	case Centered:
		return newCenteredPolyfill(resolution)
	case Intersected:
		return newIntersectedPolyfill(resolution)
	default:
		return nil
	}
}
