package db

const (
	eventsList = "list"
)

func getEventQueries() map[string]string {
	return map[string]string{
		eventsList: `
			SELECT 
				id, 
				name, 
				sport_type, 
				league, 
				location, 
				cancelled, 
				complete, 
				advertised_start_time, 
				visible 
			FROM sports_events
		`,
	}
}
