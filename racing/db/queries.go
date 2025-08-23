package db

const (
	racesList = "list"
)

func getRaceQueries() map[string]string {
	return map[string]string{
		racesList: `
			SELECT 
				id, 
				meeting_id, 
				name, 
				number, 
				visible, 
				advertised_start_time 
			FROM races
		`,
	}
}

// getSortClause returns the ORDER BY clause for the given sort option
func getSortClause(sortBy int32) string {
	defaultSort := "ORDER BY advertised_start_time ASC"
	switch sortBy {
	case 0:
		return defaultSort
	case 1:
		return "ORDER BY advertised_start_time DESC"
	case 2:
		return "ORDER BY name ASC"
	case 3:
		return "ORDER BY name DESC"
	case 4:
		return "ORDER BY number ASC"
	case 5:
		return "ORDER BY number DESC"
	default:
		return defaultSort
	}
}
