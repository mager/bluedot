package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mager/bluedot/db"
	fs "github.com/mager/bluedot/firestore"
	geojson "github.com/paulmach/go.geojson"
)

// ServeHTTP handles an HTTP requests.
func (h *Handler) getDataset(w http.ResponseWriter, r *http.Request) {
	resp := DatasetResp{}
	vars := mux.Vars(r)
	username := vars["username"]
	datasetSlug := vars["slug"]

	user := db.GetUserByUsername(h.Database, username)
	if user.ID == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	dataset := db.GetDatasetByUserIdAndSlug(h.Database, user.ID, datasetSlug)
	if dataset.ID == "" {
		http.Error(w, "Dataset not found", http.StatusNotFound)
		return
	}

	resp.ID = dataset.ID
	resp.UserID = dataset.UserID
	resp.Name = dataset.Name
	resp.Slug = dataset.Slug
	resp.Source = dataset.Source

	if dataset.Description.Valid {
		resp.Description = dataset.Description.String
	}

	if dataset.Created.Valid {
		resp.CreatedAt = dataset.Created.Time.Format("2006-01-02 15:04:05")
	}

	if dataset.Updated.Valid {
		resp.UpdatedAt = dataset.Updated.Time.Format("2006-01-02 15:04:05")
	}

	resp.User.Image = user.Image
	resp.User.Slug = user.Slug

	// Fetch record from Firestore
	doc, err := h.Firestore.Collection("datasets").Doc(resp.ID).Get(context.Background())
	if err != nil {
		h.Logger.Errorf("Error fetching Firestore record: %s", err)
		http.Error(w, "Error fetching Firestore record", http.StatusInternalServerError)
		return
	}

	img, ok := doc.Data()["image"].(string)
	if ok {
		resp.Image = img
	}

	types, ok := doc.Data()["types"].([]interface{})
	if !ok {
		resp.Types = []DatasetType{}
	} else {
		for _, t := range types {
			resp.Types = append(resp.Types, DatasetType{
				Name: fs.DatasetTypeValueToName(int(t.(int64))),
			})
		}
	}

	// Get features from Firestore
	features, err := h.Firestore.Collection("features").Where("dataset", "==", resp.ID).Documents(context.Background()).GetAll()
	if err != nil {
		h.Logger.Errorf("Error fetching features from Firestore: %s", err)
		http.Error(w, "Error fetching features from Firestore", http.StatusInternalServerError)
		return
	}
	for _, f := range features[43:44] {
		featStruct := fs.Feature{}
		err := f.DataTo(&featStruct)
		if err != nil {
			h.Logger.Errorf("Error converting Firestore data to struct: %s", err)
			http.Error(w, "Error converting Firestore data to struct", http.StatusInternalServerError)
			return
		}

		feature := geojson.Feature{}
		feature.ID = f.Ref.ID
		feature.Type = getFeatureType(featStruct.Type)
		feature.Geometry = getGeometry(featStruct)
		feature.Properties = f.Data()["properties"].(map[string]interface{})
		resp.Features = append(resp.Features, &feature)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func getFeatureType(t int) string {
	switch t {
	case 1:
		return string(geojson.GeometryPoint)
	case 2:
		return string(geojson.GeometryLineString)
	case 3:
		return string(geojson.GeometryPolygon)
	case 4:
		return string(geojson.GeometryMultiPoint)
	case 5:
		return string(geojson.GeometryMultiLineString)
	case 6:
		return string(geojson.GeometryMultiPolygon)
	case 7:
		return string(geojson.GeometryCollection)
	default:
		return "Unknown"
	}
}

func getGeometry(feat fs.Feature) *geojson.Geometry {
	geometry := &geojson.Geometry{}
	geos := feat.Geometries
	switch feat.Type {
	case fs.FeatureTypePoint:
		g := geos[0].Coords
		geometry.Type = geojson.GeometryPoint
		geometry.Geometries = append(geometry.Geometries, geojson.NewPointGeometry(g))
		geometry.Point = g
	case fs.FeatureTypeLineString:
		coords := [][]float64{}
		for _, g := range geos {
			coords = append(coords, g.Coords)
		}
		geometry.Type = geojson.GeometryLineString
		geometry.Geometries = append(geometry.Geometries, geojson.NewLineStringGeometry(coords))
		geometry.LineString = coords
	case fs.FeatureTypePolygon:
		coords := make([][][]float64, 1)
		coords[0] = make([][]float64, 0)
		for _, g := range geos {
			for i := 0; i < len(g.Coords); i += 2 {
				// Create a separate array for each pair of coordinates
				pair := []float64{g.Coords[i], g.Coords[i+1]}
				coords[0] = append(coords[0], pair)
			}
		}
		geometry.Type = geojson.GeometryPolygon
		geometry.Geometries = append(geometry.Geometries, geojson.NewPolygonGeometry(coords))
		geometry.Polygon = coords
	case fs.FeatureTypeMultiPoint:
		g := [][]float64{}
		for _, geo := range geos {
			g = append(g, geo.Coords)
		}
		geometry.Type = geojson.GeometryMultiPoint
		geometry.Geometries = append(geometry.Geometries, geojson.NewMultiPointGeometry(g...))
		geometry.MultiPoint = g
	case fs.FeatureTypeMultiLineString:
		g := [][][]float64{}
		for _, geo := range geos {
			g = append(g, [][]float64{geo.Coords})
		}
		geometry.Type = geojson.GeometryMultiLineString
		geometry.Geometries = append(geometry.Geometries, geojson.NewMultiLineStringGeometry(g...))
		geometry.MultiLineString = g
	case fs.FeatureTypeMultiPolygon:
		g := [][][][]float64{}
		for _, geo := range geos {
			g = append(g, [][][]float64{[][]float64{geo.Coords}})
		}
		geometry.Type = geojson.GeometryMultiPolygon
		geometry.Geometries = append(geometry.Geometries, geojson.NewMultiPolygonGeometry(g...))
		geometry.MultiPolygon = g
	default:
		geometry.Type = geojson.GeometryCollection
	}

	return geometry
}
