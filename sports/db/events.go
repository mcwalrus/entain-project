package db

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/sports/proto/sports"
)

// EventsRepo provides repository access to sports events.
type EventsRepo interface {
	// Init will initialise our events repository.
	Init() error

	// List will return a list of events.
	List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error)
}

type eventsRepo struct {
	db   *sql.DB
	init sync.Once
}

// NewEventsRepo creates a new events repository.
func NewEventsRepo(db *sql.DB) EventsRepo {
	return &eventsRepo{db: db}
}

func (r *eventsRepo) Init() error {
	var err error
	r.init.Do(func() {
		err = r.seed()
	})
	return err
}

func (r *eventsRepo) List(filter *sports.ListEventsRequestFilter) ([]*sports.Event, error) {
	var (
		err   error
		query string
		args  []interface{}
	)

	query = getEventQueries()[eventsList]

	query, args = r.applyFilter(query, filter)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return r.scanEvents(rows)
}

func (r *eventsRepo) applyFilter(query string, filter *sports.ListEventsRequestFilter) (string, []interface{}) {
	var (
		clauses []string
		args    []interface{}
	)

	if filter == nil {
		return query, args
	}

	if len(filter.SportTypes) > 0 {
		clauses = append(clauses, "sport_type IN ("+strings.Repeat("?,", len(filter.SportTypes)-1)+"?)")
		for _, sportType := range filter.SportTypes {
			args = append(args, sportType)
		}
	}

	if len(filter.Leagues) > 0 {
		clauses = append(clauses, "league IN ("+strings.Repeat("?,", len(filter.Leagues)-1)+"?)")
		for _, league := range filter.Leagues {
			args = append(args, league)
		}
	}

	if len(filter.Locations) > 0 {
		clauses = append(clauses, "location IN ("+strings.Repeat("?,", len(filter.Locations)-1)+"?)")
		for _, location := range filter.Locations {
			args = append(args, location)
		}
	}

	if filter.VisibleOnly {
		clauses = append(clauses, "visible = ?")
		args = append(args, 1)
	}

	if len(clauses) != 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	return query, args
}

func (r *eventsRepo) scanEvents(rows *sql.Rows) ([]*sports.Event, error) {
	var events []*sports.Event

	for rows.Next() {
		var event sports.Event
		var advertisedStart time.Time
		var cancelled, complete bool
		if err := rows.Scan(&event.Id, &event.Name, &event.SportType, &event.League, &event.Location, &cancelled, &complete, &advertisedStart, &event.Visible); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}

		ts, err := ptypes.TimestampProto(advertisedStart)
		if err != nil {
			return nil, err
		}

		event.AdvertisedStartTime = ts

		if cancelled {
			event.Status = sports.EventStatus_CANCELLED
		} else if complete {
			event.Status = sports.EventStatus_COMPLETED
		} else if advertisedStart.Before(time.Now()) {
			event.Status = sports.EventStatus_IN_PLAY
		} else {
			event.Status = sports.EventStatus_NOT_STARTED
		}

		events = append(events, &event)
	}

	return events, nil
}
