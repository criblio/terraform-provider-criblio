// Package searchacl reconciles explicit Search dataset user and team shares.
package searchacl

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/criblio/terraform-provider-criblio/internal/restclient"
)

// Grants maps principal IDs to their explicit dataset policy. A nil map leaves
// that principal type unmanaged; an empty map removes all supported shares.
type Grants struct {
	Users map[string]string `json:"users"`
	Teams map[string]string `json:"teams"`
}

// Client accesses one existing dataset in the workspace's default_search group.
type Client struct {
	rest      *restclient.Client
	datasetID string
}

// New returns a client for the specified dataset.
func New(rest *restclient.Client, datasetID string) Client {
	return Client{rest: rest, datasetID: datasetID}
}

// ValidPolicy reports whether policy is a supported dataset access level.
func ValidPolicy(policy string) bool {
	return policy == "DatasetRead" || policy == "DatasetMaintain"
}

type entry struct {
	User  string `json:"user"`
	Team  string `json:"team"`
	Perms []struct {
		Type   string `json:"type"`
		Group  string `json:"gid"`
		ID     string `json:"id"`
		Policy string `json:"policy"`
	} `json:"perms"`
}

func (c Client) path(teams bool) string {
	path := "/m/default_search/search/datasets/" + url.PathEscape(c.datasetID) + "/acl"
	if teams {
		path += "/teams"
	}
	return path
}

func (c Client) read(ctx context.Context, teams bool) (map[string]string, error) {
	items, err := restclient.Get[[]entry](ctx, c.rest, c.path(teams))
	if err != nil {
		return nil, err
	}
	if items == nil {
		return nil, fmt.Errorf("dataset ACL response contained no list")
	}
	grants := make(map[string]string)
	for _, item := range *items {
		principal := item.User
		if teams {
			principal = item.Team
		}
		for _, perm := range item.Perms {
			// Older documented responses put the dataset ID in gid and omit id.
			matches := perm.ID == c.datasetID && perm.Group == "default_search" || perm.ID == "" && perm.Group == c.datasetID
			if principal != "" && perm.Type == "datasets" && matches && ValidPolicy(perm.Policy) {
				if prior, ok := grants[principal]; ok && prior != perm.Policy {
					return nil, fmt.Errorf("principal %q has conflicting explicit dataset policies", principal)
				}
				grants[principal] = perm.Policy
			}
		}
	}
	return grants, nil
}

// Read reads selected principal types. Both reads finish before any mutation.
func (c Client) Read(ctx context.Context, users, teams bool) (Grants, error) {
	var grants Grants
	var err error
	if users {
		grants.Users, err = c.read(ctx, false)
		if err != nil {
			return Grants{}, fmt.Errorf("read user ACLs: %w", err)
		}
	}
	if teams {
		grants.Teams, err = c.read(ctx, true)
		if err != nil {
			return Grants{}, fmt.Errorf("read team ACLs: %w", err)
		}
	}
	return grants, nil
}

type changes struct {
	Add map[string][]string `json:"add,omitempty"`
	Rm  map[string][]string `json:"rm,omitempty"`
}

func (c Client) apply(ctx context.Context, desired, current map[string]string, teams bool) error {
	if desired == nil {
		return nil
	}
	body := changes{Add: make(map[string][]string), Rm: make(map[string][]string)}
	for principal, policy := range desired {
		if current[principal] != policy {
			// An add replaces the previous access level for this principal.
			body.Add[policy] = append(body.Add[policy], principal)
		}
	}
	for principal, policy := range current {
		if _, ok := desired[principal]; !ok {
			body.Rm[policy] = append(body.Rm[policy], principal)
		}
	}
	if len(body.Add) == 0 && len(body.Rm) == 0 {
		return nil
	}
	for _, policies := range []map[string][]string{body.Add, body.Rm} {
		for _, principals := range policies {
			sort.Strings(principals)
		}
	}
	return restclient.PostNoResponse(ctx, c.rest, c.path(teams)+"/apply", body)
}

// Apply reconciles configured maps with a preceding Read snapshot. The user and
// team APIs are separate transactions; a failure can follow a successful write.
func (c Client) Apply(ctx context.Context, desired, current Grants) error {
	if err := c.apply(ctx, desired.Users, current.Users, false); err != nil {
		return fmt.Errorf("apply user ACLs: %w", err)
	}
	if err := c.apply(ctx, desired.Teams, current.Teams, true); err != nil {
		return fmt.Errorf("apply team ACLs: %w", err)
	}
	return nil
}
