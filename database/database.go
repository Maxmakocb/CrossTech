package database

import (
	"fmt"
	"database/sql"
	"strconv"

	_ "github.com/lib/pq"
	
	"cross_tech/data_structures"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "asd323"
	dbname   = "cross_tech"
)

// DB is a struct, that holds a db connection session.
type DB struct{
	DB *sql.DB
}

// New creates a new DB structure and opens a connection.
func New() (*DB, error) {
	session := &DB{}
	err := session.Open()
	return session, err
}

// Open opens a new connection to a database.
func (db *DB) Open() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	var err error
	db.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return fmt.Errorf("failed, opening a db connection: %w", err)
	}

	err = db.DB.Ping()
  	if err != nil {
    	return fmt.Errorf("failed, pinging the db: %w", err)
  	}

  	fmt.Println("Successfully connected to a db")

	return nil
}

// Close closes an opened connection.
func (db *DB) Close() error {
	err := db.DB.Close()
	if err != nil {
		return fmt.Errorf("closing a db connection returned an error: %w", err)
	}
	return nil
}

func (db *DB) CreateTrack(t data_structures.Track) error {
	sqlStatement := `
	INSERT INTO track (track_id, source, target)
	VALUES ($1, $2, $3);`
	
	_, err := db.DB.Exec(sqlStatement, t.TrackId, t.Source, t.Target)
	if err != nil {
  		return fmt.Errorf("failed, adding a track with id: %d into the database: %w", t.TrackId, err)
	}

	return nil
}

func (db *DB) CreateSignal(s data_structures.SignalId, trackId int) error {
    uuid, err := strconv.Atoi(fmt.Sprintf("%d%d", trackId, s.SignalId))
    if err != nil {
        return fmt.Errorf("failed, creating a signal uuid: %w", err)
    }

	sqlStatement := `
	INSERT INTO signal_id (uuid, signal_id, track_id, signal_name, elr, mileage)
	VALUES ($1, $2, (select track_id from track where track_id = $3), $4, $5, $6);`
	
	_, err = db.DB.Exec(sqlStatement, uuid, s.SignalId, trackId, s.SignalName, s.Elr, s.Mileage)
	if err != nil {
  		return fmt.Errorf("failed, adding a signal with id: %d into the database: %w", s.SignalId, err)
	}

	return nil
}


func (db *DB) QueryTrack(trackId int) (*data_structures.Track, error) {
	response := data_structures.Track{}
	sqlStatement := `SELECT * FROM track WHERE track_id = $1;`
	row := db.DB.QueryRow(sqlStatement, trackId)
	switch err := row.Scan(&response.TrackId, &response.Source, &response.Target); err {
	case sql.ErrNoRows:
	  return nil, err
	case nil:
	default:
	  return nil, fmt.Errorf("failed, querying a track: %w", err)
	}
	return &response, nil
}

func (db *DB) QuerySignal(trackId int, signalId int) (*data_structures.SignalId, error) {
	response := data_structures.SignalId{}
	uuid, err := strconv.Atoi(fmt.Sprintf("%d%d", trackId, signalId))
    if err != nil {
        return nil, fmt.Errorf("failed, creating a signal uuid: %w", err)
    }

	sqlStatement := `SELECT * FROM signal_id WHERE uuid = $1;`
	row := db.DB.QueryRow(sqlStatement, uuid)
	switch err := row.Scan(&response.SignalId, &response.SignalId, &response.SignalName, &response.Mileage); err {
	case sql.ErrNoRows:
	  return nil, fmt.Errorf("No rows were returned!")
	case nil:
	default:
	  return nil, fmt.Errorf("failed, querying a signal: %w", err)
	}
	return &response, nil
}

func (db *DB) DeleteTrack(trackId int) error {
	sqlStatement := `DELETE FROM track WHERE track_id = $1;`
	_, err := db.DB.Exec(sqlStatement, trackId)
	if err != nil {
		return fmt.Errorf("failed, deleting a track with id: %d from the database: %w", trackId, err)
  	}
	return nil
}

func (db *DB) DeleteSignal(trackId int, signalId int) error {
	uuid, err := strconv.Atoi(fmt.Sprintf("%d%d", trackId, signalId))
    if err != nil {
        return fmt.Errorf("failed, creating a signal uuid: %w", err)
    }

	sqlStatement := `DELETE FROM signal_id WHERE uuid = $1;`
	_, err = db.DB.Exec(sqlStatement, uuid)
	if err != nil {
		return fmt.Errorf("failed, deleting a track with id: %d from the database: %w", trackId, err)
  	}
	return nil
}

var (
	TypTrackSource = "source"
	TypTrackTarget = "target"
)

func (db *DB) UpdateTrack(trackId int, newValue string, typ string) error {
	sqlStatement := ""
	switch typ {
	case TypTrackSource:
		sqlStatement = `
		UPDATE track
		SET source = $1
		WHERE track_id = $2;`
	case TypTrackTarget:
		sqlStatement = `
		UPDATE track
		SET target = $1
		WHERE track_id = $2;`
	}

	_, err := db.DB.Exec(sqlStatement, newValue, trackId)
	if err != nil {
		return fmt.Errorf("failed, updaring entry field with id: %d - %w", trackId, err)
  	}
	
	return nil
}

var (
	TypSingalMileage = "mileage"
	TypSignalElr = "elr"
	TypSignalName = "name"
)

func (db *DB) UpdateSignal(trackId int, signal_id int, newValue interface{}, typ string) error {
	uuid, err := strconv.Atoi(fmt.Sprintf("%d%d", trackId, signal_id))
    if err != nil {
        return fmt.Errorf("failed, creating a signal uuid: %w", err)
    }

	sqlStatement := ""
	switch typ {
	case TypSingalMileage:
		sqlStatement = `
		UPDATE signal_id
		SET mileage = $1
		WHERE uuid = $2;`
		update, _ := newValue.(float64)
		_, err = db.DB.Exec(sqlStatement, update, uuid)
	case TypSignalElr:
		sqlStatement = `
		UPDATE signal_id
		SET elr = $1
		WHERE uuid = $2;`
		update, _ := newValue.(string)
		_, err = db.DB.Exec(sqlStatement, update, uuid)
	case TypSignalName:
		sqlStatement = `
		UPDATE signal_id
		SET signal_name = $1
		WHERE uuid = $2;`
		update, _ := newValue.(string)
		_, err = db.DB.Exec(sqlStatement, update, uuid)
	}

	_, err = db.DB.Exec(sqlStatement, newValue, trackId)
	if err != nil {
		return fmt.Errorf("failed, updaring entry field with id: %d - %w", trackId, err)
  	}
	
	return nil
}