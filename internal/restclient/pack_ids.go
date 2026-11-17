package restclient

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const packIDCacheTTL = time.Minute

type packIDCacheEntry struct {
	ids     map[string]string
	expires time.Time
}

type packIDCache struct {
	mu      sync.Mutex
	groups  map[string]packIDCacheEntry
	lookups singleflight.Group
}

func (c *packIDCache) get(group string) (packIDCacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.groups[group]
	return entry, ok && time.Now().Before(entry.expires)
}

// ResolvePackID preserves pre-4.17 mixed-case IDs. Successful group listings are
// shared across pack resources for one minute; failures are never cached.
func (c *Client) ResolvePackID(ctx context.Context, groupID, configPackID string) string {
	lower := strings.ToLower(configPackID)
	entry, ok := c.packIDs.get(groupID)
	if !ok {
		lookup := c.packIDs.lookups.DoChan(groupID, func() (any, error) {
			if entry, ok := c.packIDs.get(groupID); ok {
				return entry, nil
			}
			items, err := List[struct {
				ID string `json:"id"`
			}](ctx, c, "/m/"+url.PathEscape(groupID)+"/packs", true)
			if err != nil {
				return nil, err
			}
			if items == nil {
				return nil, fmt.Errorf("pack list returned no response")
			}
			entry := packIDCacheEntry{ids: make(map[string]string, len(*items)), expires: time.Now().Add(packIDCacheTTL)}
			for _, item := range *items {
				entry.ids[strings.ToLower(item.ID)] = item.ID
			}
			c.packIDs.mu.Lock()
			// Bound retained group entries in long-lived debug-provider clients.
			if len(c.packIDs.groups) >= 128 {
				c.packIDs.groups = nil
			}
			if c.packIDs.groups == nil {
				c.packIDs.groups = make(map[string]packIDCacheEntry)
			}
			c.packIDs.groups[groupID] = entry
			c.packIDs.mu.Unlock()
			return entry, nil
		})
		select {
		case <-ctx.Done():
			return lower
		case result := <-lookup:
			if result.Err != nil {
				return lower
			}
			entry = result.Val.(packIDCacheEntry)
		}
	}
	if id, ok := entry.ids[lower]; ok {
		return id
	}
	return lower
}
