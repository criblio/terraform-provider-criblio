package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type datasetACLTestProvider struct{ systemLimitsTestProvider }

func (p *datasetACLTestProvider) Resources(context.Context) []func() frameworkresource.Resource {
	return []func() frameworkresource.Resource{NewSearchDatasetACLResource}
}

func (p *datasetACLTestProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{NewSearchDatasetACLDataSource}
}

type datasetACLBackend struct {
	mu           sync.Mutex
	t            *testing.T
	users, teams map[string]string
	posts        int
	failRead     bool
	failTeamOnce bool
	missing      bool
}

func (b *datasetACLBackend) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	b.mu.Lock()
	defer b.mu.Unlock()
	const base = "/api/v1/m/default_search/search/datasets/logs/acl"
	if !strings.HasPrefix(req.URL.Path, base) || b.missing {
		http.NotFound(w, req)
		return
	}
	teams := strings.Contains(req.URL.Path, "/teams")
	stored := b.users
	principalKey := "user"
	if teams {
		stored, principalKey = b.teams, "team"
	}
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodGet:
		if b.failRead && teams {
			http.Error(w, "read denied", http.StatusForbidden)
			return
		}
		items := []any{}
		for principal, policy := range stored {
			items = append(items, map[string]any{principalKey: principal, "perms": []any{
				map[string]string{"type": "datasets", "gid": "default_search", "id": "logs", "policy": policy},
				map[string]string{"type": "datasets", "gid": "default_search", "id": "other", "policy": "DatasetRead"},
			}})
		}
		items = append(items, map[string]any{principalKey: "inherited", "perms": []any{
			map[string]string{"type": "groups", "gid": "default_search", "policy": "GroupFull"},
		}})
		if err := json.NewEncoder(w).Encode(map[string]any{"items": items, "count": len(items)}); err != nil {
			b.t.Errorf("encode ACL list: %v", err)
		}
	case http.MethodPost:
		b.posts++
		if req.URL.Path != base+"/apply" && req.URL.Path != base+"/teams/apply" {
			b.t.Errorf("unexpected ACL action: %s", req.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if teams && b.failTeamOnce {
			b.failTeamOnce = false
			http.Error(w, "team apply denied", http.StatusForbidden)
			return
		}
		var body map[string]map[string][]string
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			b.t.Errorf("decode action: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(body) == 0 {
			b.t.Error("sent an empty ACL action")
		}
		for action := range body {
			if action != "add" && action != "rm" {
				b.t.Errorf("declarative field leaked into action: %s", action)
			}
		}
		for policy, principals := range body["rm"] {
			for _, principal := range principals {
				if stored[principal] == policy {
					delete(stored, principal)
				}
			}
		}
		for policy, principals := range body["add"] {
			for _, principal := range principals {
				stored[principal] = policy
			}
		}
		if _, err := w.Write([]byte(`{"items":[],"count":0}`)); err != nil {
			b.t.Errorf("write ACL action response: %v", err)
		}
	default:
		b.t.Errorf("unexpected ACL method: %s", req.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func datasetACLFactories(url string) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"criblio": providerserver.NewProtocol6WithError(&datasetACLTestProvider{systemLimitsTestProvider{baseURL: url}}),
	}
}

func aclTestConfig(maps string) string {
	return fmt.Sprintf(`resource "criblio_search_dataset_acl" "test" {
  dataset_id = "logs"
  %s
}`, maps)
}

func TestSearchDatasetACLLifecycle(t *testing.T) {
	backend := &datasetACLBackend{t: t, users: map[string]string{"existing@example.com": "DatasetMaintain", "custom": "FutureDatasetPolicy"}, teams: map[string]string{"old-team": "DatasetRead"}}
	server := httptest.NewServer(backend)
	defer server.Close()
	name := "criblio_search_dataset_acl.test"
	teamConfig := aclTestConfig(`teams = { engineering = "DatasetRead" }`)
	allConfig := aclTestConfig(`teams = { engineering = "DatasetMaintain" }
users = { "google-oauth2|example-member-id" = "DatasetRead", "client-id@clients" = "DatasetRead" }`)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: datasetACLFactories(server.URL),
		Steps: []resource.TestStep{
			{Config: teamConfig, Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "teams.engineering", "DatasetRead"),
				resource.TestCheckNoResourceAttr(name, "users"),
			)},
			{Config: teamConfig, PlanOnly: true},
			{Config: allConfig + "\ndata \"criblio_search_dataset_acl\" \"test\" {\n  dataset_id = criblio_search_dataset_acl.test.dataset_id\n  depends_on = [criblio_search_dataset_acl.test]\n}", Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttr(name, "users.client-id@clients", "DatasetRead"),
				resource.TestCheckResourceAttr("data.criblio_search_dataset_acl.test", "teams.engineering", "DatasetMaintain"),
			)},
			{ResourceName: name, ImportState: true, ImportStateId: "logs", ImportStateVerify: true, ImportStateVerifyIdentifierAttribute: "dataset_id"},
			{PreConfig: func() {
				backend.mu.Lock()
				defer backend.mu.Unlock()
				backend.teams["engineering"] = "DatasetRead"
				backend.users["drift@example.com"] = "DatasetRead"
			}, Config: allConfig, PlanOnly: true, ExpectNonEmptyPlan: true},
			{Config: allConfig},
			{Config: allConfig, PlanOnly: true},
			{Config: aclTestConfig(`users = {}`)},
			{Config: aclTestConfig(`users = {}`), PlanOnly: true},
		},
	})
	backend.mu.Lock()
	defer backend.mu.Unlock()
	if !reflect.DeepEqual(backend.users, map[string]string{"custom": "FutureDatasetPolicy"}) || !reflect.DeepEqual(backend.teams, map[string]string{"engineering": "DatasetMaintain"}) {
		t.Fatalf("destroy or reconciliation affected unmanaged shares: users=%v teams=%v", backend.users, backend.teams)
	}
	if backend.posts != 6 {
		t.Errorf("action count = %d, want 6 (no redundant writes)", backend.posts)
	}
}

func TestSearchDatasetACLPartialFailureRecovery(t *testing.T) {
	backend := &datasetACLBackend{t: t, users: map[string]string{}, teams: map[string]string{}, failTeamOnce: true}
	server := httptest.NewServer(backend)
	defer server.Close()
	config := aclTestConfig(`users = { "client-id@clients" = "DatasetRead" }
teams = { engineering = "DatasetRead" }`)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: datasetACLFactories(server.URL),
		Steps: []resource.TestStep{
			{Config: aclTestConfig("users = {}\nteams = {}")},
			{Config: config, ExpectError: regexp.MustCompile("Unable to reconcile dataset ACLs")},
			{Config: config},
			{Config: config, PlanOnly: true},
		},
	})
	if len(backend.users) != 0 || len(backend.teams) != 0 {
		t.Fatal("destroy did not remove managed shares")
	}
	if backend.posts != 5 {
		t.Errorf("action count = %d, want 5 (retry must not reapply user grants)", backend.posts)
	}
}

func TestSearchDatasetACLReadFailurePreventsActions(t *testing.T) {
	backend := &datasetACLBackend{t: t, users: map[string]string{}, teams: map[string]string{}, failRead: true}
	server := httptest.NewServer(backend)
	defer server.Close()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: datasetACLFactories(server.URL),
		Steps: []resource.TestStep{{Config: aclTestConfig(`users = { "client-id@clients" = "DatasetRead" }
teams = { engineering = "DatasetRead" }`), ExpectError: regexp.MustCompile("read team ACLs")}},
	})
	if backend.posts != 0 {
		t.Fatal("a failed preflight read caused an ACL write")
	}
}

func TestSearchDatasetACLInvalidPolicies(t *testing.T) {
	for _, maps := range []string{`teams = { engineering = "Admin" }`, `users = { "client-id@clients" = null }`} {
		resource.UnitTest(t, resource.TestCase{
			ProtoV6ProviderFactories: datasetACLFactories("http://127.0.0.1:1"),
			Steps:                    []resource.TestStep{{Config: aclTestConfig(maps), PlanOnly: true, ExpectError: regexp.MustCompile("Invalid dataset policy")}},
		})
	}
}

func TestSearchDatasetACLUnknownReferences(t *testing.T) {
	backend := &datasetACLBackend{t: t, users: map[string]string{}, teams: map[string]string{}}
	server := httptest.NewServer(backend)
	defer server.Close()
	config := `resource "terraform_data" "dataset" { input = "logs" }
resource "terraform_data" "principal" { input = "client-id@clients" }
resource "criblio_search_dataset_acl" "test" {
  dataset_id = terraform_data.dataset.output
  users = { (terraform_data.principal.output) = "DatasetRead" }
}`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: datasetACLFactories(server.URL),
		Steps: []resource.TestStep{
			{Config: config, Check: resource.TestCheckResourceAttr("criblio_search_dataset_acl.test", "users.client-id@clients", "DatasetRead")},
			{Config: config, PlanOnly: true},
		},
	})
}

func TestSearchDatasetACLDeletedDataset(t *testing.T) {
	backend := &datasetACLBackend{t: t, users: map[string]string{}, teams: map[string]string{}}
	server := httptest.NewServer(backend)
	defer server.Close()
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: datasetACLFactories(server.URL),
		Steps: []resource.TestStep{
			{Config: aclTestConfig(`teams = { engineering = "DatasetRead" }`)},
			{PreConfig: func() {
				backend.mu.Lock()
				defer backend.mu.Unlock()
				backend.missing = true
			}, Config: `terraform {}`},
		},
	})
	if backend.posts != 1 {
		t.Fatal("deleting a missing dataset attempted an ACL mutation")
	}
}
