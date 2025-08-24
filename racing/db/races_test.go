package db

import (
	"database/sql"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/proto/racing"
)

var (
	testRepo RacesRepo
)

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite3", "./racing.db")
	if err != nil {
		panic("Failed to open racing.db: " + err.Error())
	}

	testRepo = NewRacesRepo(db)
	if err := testRepo.Init(); err != nil {
		panic("Failed to initialize races repository: " + err.Error())
	}

	allRaces, err := testRepo.List(nil)
	if err != nil {
		panic("Failed to get all races: " + err.Error())
	}

	if len(allRaces) == 0 {
		panic("No races available for testing")
	}

	code := m.Run()
	os.Exit(code)
}

func contains(ids []int64, id int64) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

func TestRacesRepo_List_MeetingIds(t *testing.T) {

	t.Run("with no filter", func(t *testing.T) {
		filter := &racing.ListRacesRequestFilter{
			MeetingIds: []int64{},
		}

		races, err := testRepo.List(filter)
		if err != nil {
			t.Fatalf("List() with empty meeting_ids filter failed: %v", err)
		}

		if len(races) == 0 {
			t.Fatal("Expected some races, got none")
		}
	})

	t.Run("with meeting_ids filter", func(t *testing.T) {
		testMeetingIDs := []int64{1, 2, 3}
		filter := &racing.ListRacesRequestFilter{
			MeetingIds: testMeetingIDs,
		}

		races, err := testRepo.List(filter)
		if err != nil {
			t.Fatalf("List() with meeting_ids filter failed: %v", err)
		}

		if len(races) == 0 {
			t.Fatal("Expected some races, got none")
		}

		for _, race := range races {
			if !contains(testMeetingIDs, race.MeetingId) {
				t.Errorf("Expected race with meeting ID %d, got %d", testMeetingIDs, race.MeetingId)
			}
		}
	})
}

func TestRacesRepo_List_VisibleOnlyFilter(t *testing.T) {

	t.Run("expect both visible and invisible races when visible_only=false", func(t *testing.T) {
		filter := &racing.ListRacesRequestFilter{
			VisibleOnly: false,
		}

		races, err := testRepo.List(filter)
		if err != nil {
			t.Fatalf("List() with visible_only filter failed: %v", err)
		}

		visibleCount := 0
		invisibleCount := 0
		for _, race := range races {
			if race.Visible {
				visibleCount++
			} else {
				invisibleCount++
			}
		}

		if visibleCount == 0 || invisibleCount == 0 {
			t.Fatalf("Expected some visible and invisible races, got %d visible and %d invisible", visibleCount, invisibleCount)
		}

	})

	t.Run("expect only visible races when visible_only=true", func(t *testing.T) {
		filter := &racing.ListRacesRequestFilter{
			VisibleOnly: true,
		}

		races, err := testRepo.List(filter)
		if err != nil {
			t.Fatalf("List() with visible_only filter failed: %v", err)
		}

		for _, race := range races {
			if !race.Visible {
				t.Errorf("Expected only visible races, but got race ID %d with visible=%t", race.Id, race.Visible)
			}
		}
	})
}

func verifySort(t *testing.T, race1, race2 *racing.Race, sortBy racing.ListRacesSortBy) {
	t.Helper()

	switch sortBy {
	case racing.ListRacesSortBy_ADVERTISED_START_TIME_ASC:
		time1, err := ptypes.Timestamp(race1.AdvertisedStartTime)
		if err != nil {
			t.Fatalf("Failed to get advertised_start_time for race %d: %v", race1.Id, err)
		}
		time2, err := ptypes.Timestamp(race2.AdvertisedStartTime)
		if err != nil {
			t.Fatalf("Failed to get advertised_start_time for race %d: %v", race2.Id, err)
		}
		if time1.After(time2) {
			t.Error("Races not sorted by advertised_start_time ASC")
		}
	case racing.ListRacesSortBy_ADVERTISED_START_TIME_DESC:
		time1, err := ptypes.Timestamp(race1.AdvertisedStartTime)
		if err != nil {
			t.Fatalf("Failed to get advertised_start_time for race %d: %v", race1.Id, err)
		}
		time2, err := ptypes.Timestamp(race2.AdvertisedStartTime)
		if err != nil {
			t.Fatalf("Failed to get advertised_start_time for race %d: %v", race2.Id, err)
		}
		if time1.Before(time2) {
			t.Error("Races not sorted by advertised_start_time DESC")
		}
	case racing.ListRacesSortBy_NAME_ASC:
		if race1.Name > race2.Name {
			t.Error("Races not sorted by name ASC")
		}
	case racing.ListRacesSortBy_NAME_DESC:
		if race1.Name < race2.Name {
			t.Error("Races not sorted by name DESC")
		}
	case racing.ListRacesSortBy_NUMBER_ASC:
		if race1.Number > race2.Number {
			t.Error("Races not sorted by number ASC")
		}
	case racing.ListRacesSortBy_NUMBER_DESC:
		if race1.Number < race2.Number {
			t.Error("Races not sorted by number DESC")
		}
	}
}

func TestRacesRepo_List_SortOptions(t *testing.T) {
	tests := []struct {
		name   string
		sortBy racing.ListRacesSortBy
	}{
		{
			name:   "ADVERTISED_START_TIME_ASC",
			sortBy: racing.ListRacesSortBy_ADVERTISED_START_TIME_ASC,
		},
		{
			name:   "ADVERTISED_START_TIME_DESC",
			sortBy: racing.ListRacesSortBy_ADVERTISED_START_TIME_DESC,
		},
		{
			name:   "NAME_ASC",
			sortBy: racing.ListRacesSortBy_NAME_ASC,
		},
		{
			name:   "NAME_DESC",
			sortBy: racing.ListRacesSortBy_NAME_DESC,
		},
		{
			name:   "NUMBER_ASC",
			sortBy: racing.ListRacesSortBy_NUMBER_ASC,
		},
		{
			name:   "NUMBER_DESC",
			sortBy: racing.ListRacesSortBy_NUMBER_DESC,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := &racing.ListRacesRequestFilter{
				SortBy: tt.sortBy,
			}

			races, err := testRepo.List(filter)
			if err != nil {
				t.Fatalf("List() with sort %v failed: %v", tt.sortBy, err)
			}

			if len(races) == 0 {
				t.Fatal("Expected races, got none")
			}

			if len(races) < 2 {
				return
			}

			for i := 0; i < len(races)-1; i++ {
				verifySort(races[i], races[i+1])
			}
		})
	}
}

func TestRacesRepo_List_SortWithFilters(t *testing.T) {
	filter := &racing.ListRacesRequestFilter{
		VisibleOnly: true,
		SortBy:      racing.ListRacesSortBy_NAME_ASC,
	}

	races, err := testRepo.List(filter)
	if err != nil {
		t.Fatalf("List() with visible filter and sort failed: %v", err)
	}

	for _, race := range races {
		if !race.Visible {
			t.Errorf("Race should be visible")
		}
	}

	if len(races) >= 2 && races[0].Name > races[1].Name {
		t.Error("Races not sorted by name ASC")
	}
}
