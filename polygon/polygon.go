package polygon

import (
	geojson "github.com/paulmach/go.geojson"
	"github.com/tidwall/geojson/geometry"
	"github.com/uber/h3-go/v3"
	"strconv"
)

type polygon [][]float64

func Intersects(a, b geojson.Feature) bool {
	return polyToGeometry(a.Geometry.Polygon[0]).IntersectsPoly(polyToGeometry(b.Geometry.Polygon[0]))
}

func FeatureToH3GeoPolygon(feature geojson.Feature) h3.GeoPolygon {
	geoFence := make([]h3.GeoCoord, len(feature.Geometry.Polygon[0]))
	for i, coords := range feature.Geometry.Polygon[0] {
		geoFence[i] = h3.GeoCoord{
			Longitude: coords[0],
			Latitude:  coords[1],
		}
	}

	return h3.GeoPolygon{
		Geofence: geoFence,
		Holes:    nil,
	}
}

func HexagonToFeature(hexagon uint64) geojson.Feature {
	var polygon polygon
	hexBoundaries := h3.ToGeoBoundary(h3.FromString(strconv.FormatUint(hexagon, 16)))
	for _, boundary := range hexBoundaries {
		polygon = append(polygon, []float64{boundary.Longitude, boundary.Latitude})
	}

	geoJsonPolygon := append(reverse(polygon), polygon[len(polygon)-1])
	feature := geojson.NewFeature(geojson.NewPolygonGeometry([][][]float64{geoJsonPolygon}))
	return *feature
}

func reverse(s [][]float64) [][]float64 {
	a := make([][]float64, len(s))
	copy(a, s)

	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}

	return a
}

func polyToGeometry(p polygon) *geometry.Poly {
	holes := make([][]geometry.Point, 0)
	points := make([]geometry.Point, 0)
	for _, point := range p {
		points = append(points, geometry.Point{
			X: point[0],
			Y: point[1],
		})
	}

	return geometry.NewPoly(points, holes, &geometry.IndexOptions{})
}
