package polyfill

import (
	"github.com/jeromeshapiro/geoToH3/polygon"
	geojson "github.com/paulmach/go.geojson"
	"github.com/uber/h3-go/v3"
)

type CenteredPolyfill struct {
	resolution int
}

func newCenteredPolyfill(resolution int) CenteredPolyfill {
	return CenteredPolyfill{resolution}
}

func (c CenteredPolyfill) featureToH3(feature geojson.Feature) []uint64 {
	hexagons := make([]uint64, 0)

	polyfill := h3.Polyfill(polygon.FeatureToH3GeoPolygon(feature), c.resolution)

	for _, h3Index := range polyfill {
		hexagons = append(hexagons, uint64(h3Index))
	}

	return hexagons
}

func (c CenteredPolyfill) FeatureCollectionToH3(featureCollection geojson.FeatureCollection) []uint64 {
	hexagons := make([]uint64, 0)

	for _, feature := range featureCollection.Features {
		hexagons = append(hexagons, c.featureToH3(*feature)...)
	}

	return hexagons
}
