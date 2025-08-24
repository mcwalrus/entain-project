package db

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/mattn/go-sqlite3"

	"git.neds.sh/matty/entain/racing/pkg/clock"
	"git.neds.sh/matty/entain/racing/proto/racing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		require.NoError(t, err, "List() with empty meeting_ids filter failed")
		assert.NotEmpty(t, races, "Expected some races, got none")
	})

	t.Run("with meeting_ids filter", func(t *testing.T) {
		testMeetingIDs := []int64{1, 2, 3}
		filter := &racing.ListRacesRequestFilter{
			MeetingIds: testMeetingIDs,
		}

		races, err := testRepo.List(filter)
		require.NoError(t, err, "List() with meeting_ids filter failed")
		assert.NotEmpty(t, races, "Expected some races, got none")

		for _, race := range races {
			assert.True(t, contains(testMeetingIDs, race.MeetingId), "Expected race with meeting ID in %v, got %d", testMeetingIDs, race.MeetingId)
		}
	})
}

func TestRacesRepo_List_VisibleOnlyFilter(t *testing.T) {

	t.Run("expect both visible and invisible races when visible_only=false", func(t *testing.T) {
		filter := &racing.ListRacesRequestFilter{
			VisibleOnly: false,
		}

		races, err := testRepo.List(filter)
		require.NoError(t, err, "List() with visible_only filter failed")

		visibleCount := 0
		invisibleCount := 0
		for _, race := range races {
			if race.Visible {
				visibleCount++
			} else {
				invisibleCount++
			}
		}

		assert.Greater(t, visibleCount, 0, "Expected some visible races")
		assert.Greater(t, invisibleCount, 0, "Expected some invisible races")

	})

	t.Run("expect only visible races when visible_only=true", func(t *testing.T) {
		filter := &racing.ListRacesRequestFilter{
			VisibleOnly: true,
		}

		races, err := testRepo.List(filter)
		require.NoError(t, err, "List() with visible_only filter failed")

		for _, race := range races {
			assert.True(t, race.Visible, "Expected only visible races, but got race ID %d with visible=%t", race.Id, race.Visible)
		}
	})
}

func verifySort(t *testing.T, race1, race2 *racing.Race, sortBy racing.ListRacesSortBy) {
	t.Helper()

	time1, err := ptypes.Timestamp(race1.AdvertisedStartTime)
	require.NoError(t, err)
	time2, err := ptypes.Timestamp(race2.AdvertisedStartTime)
	require.NoError(t, err)

	switch sortBy {
	case racing.ListRacesSortBy_ADVERTISED_START_TIME_ASC:
		assert.LessOrEqual(t, time1, time2)
	case racing.ListRacesSortBy_ADVERTISED_START_TIME_DESC:
		assert.GreaterOrEqual(t, time1, time2)
	case racing.ListRacesSortBy_NAME_ASC:
		assert.LessOrEqual(t, race1.Name, race2.Name)
	case racing.ListRacesSortBy_NAME_DESC:
		assert.GreaterOrEqual(t, race1.Name, race2.Name)
	case racing.ListRacesSortBy_NUMBER_ASC:
		assert.LessOrEqual(t, race1.Number, race2.Number)
	case racing.ListRacesSortBy_NUMBER_DESC:
		assert.GreaterOrEqual(t, race1.Number, race2.Number)
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
			require.NoError(t, err, "List() with sort %v failed", tt.sortBy)
			require.NotEmpty(t, races, "Expected races, got none")
			require.GreaterOrEqual(t, len(races), 2, "Expected at least 2 races, got %d", len(races))

			for i := 0; i < len(races)-1; i++ {
				verifySort(t, races[i], races[i+1], tt.sortBy)
			}
		})
	}
}

func TestScanRaces_RaceStatus(t *testing.T) {
	t.Cleanup(func() {
		clock.ResetClockImplementation()
	})

	filter := &racing.ListRacesRequestFilter{
		SortBy: racing.ListRacesSortBy_ADVERTISED_START_TIME_ASC,
	}

	races, err := testRepo.List(filter)
	require.NoError(t, err)
	require.NotEmpty(t, races)

	testTime := races[0].AdvertisedStartTime.AsTime()

	t.Run("races before test time should be OPEN", func(t *testing.T) {
		clock.NowFunc(func() time.Time { return testTime.Add(+1 * time.Second) })

		races, err := testRepo.List(filter)
		require.NoError(t, err)

		require.Less(t, races[0].AdvertisedStartTime.AsTime(), clock.Now())
		require.Equal(t, racing.RaceStatus_OPEN, races[0].Status)
	})

	t.Run("races at test time should be CLOSED", func(t *testing.T) {
		clock.NowFunc(func() time.Time { return testTime })

		races, err := testRepo.List(filter)
		require.NoError(t, err)

		require.Equal(t, races[0].AdvertisedStartTime.AsTime(), clock.Now())
		require.Equal(t, racing.RaceStatus_CLOSED, races[0].Status)
	})

	t.Run("races after test time should be CLOSED", func(t *testing.T) {
		clock.NowFunc(func() time.Time { return testTime.Add(-1 * time.Second) })

		races, err := testRepo.List(filter)
		require.NoError(t, err)

		require.Greater(t, races[0].AdvertisedStartTime.AsTime(), clock.Now())
		require.Equal(t, racing.RaceStatus_CLOSED, races[0].Status)
	})
}

func TestRacesRepo_Get(t *testing.T) {
	allRaces, err := testRepo.List(nil)
	require.NoError(t, err)
	require.NotEmpty(t, allRaces)

	t.Run("get existing race by ID", func(t *testing.T) {
		expectedRace := allRaces[0]
		race, err := testRepo.Get(expectedRace.Id)
		require.NoError(t, err)
		require.NotNil(t, race)
	})

	t.Run("get non-existent race by ID", func(t *testing.T) {
		nonExistentID := int64(99999)
		race, err := testRepo.Get(nonExistentID)
		require.NoError(t, err)
		assert.Nil(t, race)
	})
}
