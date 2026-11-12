package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
)

func TestDatasetACLSharedTeamFixture(t *testing.T) {
	for _, tc := range []struct {
		name       string
		exists     bool
		lookupCode int
		createCode int
		wantError  bool
		wantCalls  []string
	}{
		{name: "reuse existing team", exists: true, wantCalls: []string{"GET"}},
		{name: "create missing team", wantCalls: []string{"GET", "POST", "GET"}},
		{name: "concurrent creation", createCode: http.StatusConflict, wantCalls: []string{"GET", "POST", "GET"}},
		{name: "lookup denied", lookupCode: http.StatusForbidden, wantError: true, wantCalls: []string{"GET"}},
		{name: "creation denied", createCode: http.StatusForbidden, wantError: true, wantCalls: []string{"GET", "POST", "GET"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			exists := tc.exists
			var calls []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				calls = append(calls, req.Method)
				w.Header().Set("Content-Type", "application/json")
				switch req.Method {
				case http.MethodGet:
					if req.URL.Path != "/api/v1/system/teams/test-team" {
						t.Errorf("unexpected lookup path: %s", req.URL.Path)
					}
					if tc.lookupCode != 0 {
						w.WriteHeader(tc.lookupCode)
						return
					}
					if !exists {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					if err := json.NewEncoder(w).Encode(map[string]any{"items": []datasetACLTestTeam{{ID: "test-team"}}}); err != nil {
						t.Error(err)
					}
				case http.MethodPost:
					if req.URL.Path != "/api/v1/system/teams" {
						t.Errorf("unexpected creation path: %s", req.URL.Path)
					}
					var team datasetACLTestTeam
					if err := json.NewDecoder(req.Body).Decode(&team); err != nil {
						t.Error(err)
					}
					if team.ID != "test-team" || team.Name != "test-team" || !reflect.DeepEqual(team.Roles, []string{"search_user"}) {
						t.Errorf("unexpected fixture: %+v", team)
					}
					exists = tc.createCode == 0 || tc.createCode == http.StatusConflict
					if tc.createCode != 0 {
						w.WriteHeader(tc.createCode)
						return
					}
					w.WriteHeader(http.StatusNoContent)
				default:
					t.Errorf("shared team must not be modified or deleted: %s", req.Method)
					w.WriteHeader(http.StatusMethodNotAllowed)
				}
			}))
			defer server.Close()
			client := restclient.New(restclient.Config{BaseURL: server.URL, BearerToken: "test"})
			err := ensureDatasetACLTestTeam(t.Context(), client)
			if (err != nil) != tc.wantError {
				t.Fatalf("fixture error = %v, want error: %v", err, tc.wantError)
			}
			if !reflect.DeepEqual(calls, tc.wantCalls) {
				t.Fatalf("fixture calls = %v, want %v", calls, tc.wantCalls)
			}
		})
	}
}
