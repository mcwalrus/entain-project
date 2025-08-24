package db

import (
	"time"

	"syreclabs.com/go/faker"
)

func (r *eventsRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS sports_events (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		sport_type TEXT NOT NULL,
		league TEXT,
		location TEXT,
		cancelled BOOLEAN,
		complete BOOLEAN,
		advertised_start_time DATETIME NOT NULL,
		visible INTEGER DEFAULT 1
	)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO sports_events(id, name, sport_type, league, location, cancelled, complete, advertised_start_time, visible) VALUES (?,?,?,?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Team().Name()+" vs "+faker.Team().Name(),
				faker.RandomChoice([]string{"Football", "Basketball", "Tennis", "Baseball", "Soccer", "Hockey"}),
				faker.Company().Name()+" League",
				faker.Address().City(),
				faker.Number().Between(0, 1) == "1",
				faker.Number().Between(0, 1) == "1",
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
				faker.Number().Between(0, 1),
			)
		}
	}

	return err
}
