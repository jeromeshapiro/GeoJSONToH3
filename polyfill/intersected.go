package polyfill

import (
	"container/list"
	"github.com/jeromeshapiro/geoToH3/polygon"
	geojson "github.com/paulmach/go.geojson"
	"github.com/uber/h3-go/v3"
)

type IntersectedPolyfill struct {
	resolution int
}

func newIntersectedPolyfill(resolution int) IntersectedPolyfill {
	return IntersectedPolyfill{resolution}
}

func (i IntersectedPolyfill) FeatureToH3(feature geojson.Feature) []uint64 {
	found := make([]h3.H3Index, 0)
	visited := make(map[h3.H3Index]struct{})

	queue := list.New()

	for _, hexagon := range h3.Polyfill(polygon.FeatureToH3GeoPolygon(feature), i.resolution) {
		visited[hexagon] = struct{}{}
		queue.PushBack(hexagon)
	}

	for queue.Len() > 0 {
		node := queue.Front()
		hexagon := node.Value.(h3.H3Index)
		queue.Remove(node)

		found = append(found, hexagon)

		for _, neighbour := range h3.KRing(hexagon, 1) {
			if _, hasVisited := visited[neighbour]; !hasVisited {
				visited[neighbour] = struct{}{}

				if polygon.Intersects(feature, polygon.HexagonToFeature(uint64(neighbour))) {
					queue.PushBack(neighbour)
				}
			}
		}
	}

	var hexInts []uint64
	for _, h := range found {
		hexInts = append(hexInts, uint64(h))
	}

	return hexInts
}

func (i IntersectedPolyfill) FeatureCollectionToH3(featureCollection geojson.FeatureCollection) []uint64 {
	hexagons := make([]uint64, 0)

	for _, feature := range featureCollection.Features {
		hexagons = append(hexagons, i.FeatureToH3(*feature)...)
	}

	return hexagons
}
