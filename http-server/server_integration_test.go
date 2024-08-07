package poker

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	type testConfig struct {
		name  string
		store PlayerStore
		want  []Player
	}
	postgresPlayerStore := NewPostgresPlayerStore()
	t.Cleanup(postgresPlayerStore.Close)
	players := []Player{
		{"Pepper", 3},
	}
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)

	for _, tc := range []testConfig{
		{
			name:  "InMemoryPlayerStore",
			store: store,
			want:  players,
		},
		//{
		//	name:  "PostgresPlayerStore",
		//	store: postgresPlayerStore,
		//	want:  players,
		//},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := NewPlayerServer(tc.store)
			player := "Pepper"

			server.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil))
			server.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil))
			server.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", player), nil))

			response := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", player), nil)
			server.ServeHTTP(response, request)

			AssertEqual(t, response.Code, http.StatusOK)
			AssertEqual(t, response.Body.String(), "3")

			request, _ = http.NewRequest(http.MethodGet, "/league", nil)
			response = httptest.NewRecorder()
			server.ServeHTTP(response, request)
			AssertEqual(t, response.Code, http.StatusOK)

			got := GetLeagueFromResponse(t, response)
			AssertLeague(t, got, tc.want)
		})
	}
}
