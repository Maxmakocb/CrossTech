package server

import (
	"fmt"
	"database/sql"
	"net/http"
	"encoding/json"
	"strconv"

	"github.com/labstack/echo"

	"cross_tech/data_structures"
	"cross_tech/database"
)

type srv struct {
	server *echo.Echo
	db     *database.DB
}

// New creates an instance of the server.
func New() (*srv, error) {
	instance := srv{
		server: echo.New(),
	}

	db, err := database.New()
	if err != nil {
		return nil, err
	}

	instance.db = db

	instance.server.GET("track", func(c echo.Context) error {
		trackId := c.QueryParam("id")

		fmt.Printf("GET request, query for track with id: %s\n", trackId)

		id, err := strconv.Atoi(trackId)
    	if err != nil {
			fmt.Printf("GET request fail, bad request: %v\n", err)
			return c.String(http.StatusBadRequest, fmt.Sprintf("track id was not an integer: %s, error: %v", trackId, err))
    	}

		track, err := db.QueryTrack(id)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.String(http.StatusOK, "No entries were found")
			}
			fmt.Printf("GET request fail, track not found: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed, querying a track with id: %d, error: %v", id, err))
		}
		fmt.Printf("GET request success, track found: %v\n", track)
		return c.String(http.StatusOK, MarshalDatastruct(*track))
	})

	instance.server.DELETE("track", func(c echo.Context) error {
		trackId := c.QueryParam("id")

		fmt.Printf("DELETE request for track with id: %s\n", trackId)

		id, err := strconv.Atoi(trackId)
    	if err != nil {
			fmt.Printf("DELETE request fail, bad request: %v\n", err)
			return c.String(http.StatusBadRequest, fmt.Sprintf("track id was not an integer: %s, error: %v", trackId, err))
    	}

		err = db.DeleteTrack(id)
		if err != nil {
			fmt.Printf("DELETE request fail: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed, deleting a track with id: %d, error: %v", id, err))
		}
		fmt.Printf("DELETE request success, track deleted\n")
		return c.String(http.StatusOK, fmt.Sprintf("Track with id %d has been deleted", trackId))
	})

	instance.server.PUT("track", func(c echo.Context) error {
		trackId := c.QueryParam("id")
		typ := c.QueryParam("typ")
		value := c.QueryParam("value")

		fmt.Printf("PUT request for track with id: %s\n", trackId)

		id, err := strconv.Atoi(trackId)
    	if err != nil {
			fmt.Printf("PUT request fail, bad request: %v\n", err)
			return c.String(http.StatusBadRequest, fmt.Sprintf("track id was not an integer: %s, error: %v", trackId, err))
    	}

		_, err = db.QueryTrack(id)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.String(http.StatusOK, "No entries were found")
			}
			fmt.Printf("PUT request fail, track not found: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed, querying a track with id: %d, error: %v", id, err))
		}

		err = db.UpdateTrack(id, value, typ) 
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed, updating a track with id: %d, error: %v", id, err))
		}

		fmt.Printf("PUT request success, track update: %v\n", err)
		return c.String(http.StatusOK, "successfully updated a track")
	})

	instance.server.POST("/track", func(c echo.Context) error {
		jsonEntry := c.QueryParam("entry")
		fmt.Printf("POST request for a track with input json: %v\n", jsonEntry)
		var entry data_structures.TrackStandalone
		err := json.Unmarshal([]byte(jsonEntry), &entry)
		if err != nil {
			fmt.Printf("POST request fail, failed, parsing json argument: %v\n", err)
			return c.String(http.StatusBadRequest, fmt.Sprintf("failed, parsing json argument: %v", err))
		}

		track, err := db.QueryTrack(entry.TrackId)
		if track != nil {
			return c.String(http.StatusOK, "Track already exists, cannot add a new track in place")
		}
		if err != nil && sql.ErrNoRows != err {
			fmt.Printf("POST request fail, track not found: %v\n", err)
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed, querying a track with id: %d, error: %v", entry.TrackId, err))
		}

		err = db.CreateTrack(data_structures.Track{
			TrackId: entry.TrackId,
			Source: entry.Source,
			Target: entry.Target,
		}) 
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed, creating a track with id: %d, error: %v", entry.TrackId, err))
		}

		fmt.Printf("POST request success, track update: %v\n", err)
		return c.String(http.StatusOK, "successfully created a track")
	})

	// TODO
	// Add the methods to manage signal_id data.

	return &instance, nil
}

func (s *srv) Start(port int) error {
	return s.server.Start(fmt.Sprintf(":%d", port))
}

func MarshalDatastruct(datastruct interface{}) string {
	if track, ok := datastruct.(data_structures.Track); ok {
		standalone := data_structures.TrackToStandalone(track)
		output, _ := json.Marshal(standalone)
		return string(output)
	} else if signal, ok := datastruct.(data_structures.SignalId); ok {
		output, _ := json.Marshal(signal)
		return string(output)
	}
	return ""
}