package restclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	listPageSize = 100
	// Bound requests even when an API ignores pagination and returns volatile
	// pages that cannot be recognized by byte equality.
	listMaxPages = 1000
)

// List reads a collection, following Cribl offset/limit pagination when supported.
// APIs without pagination continue to use their original request path.
func List[T any](ctx context.Context, c *Client, path string, paginated bool) (*[]T, error) {
	requestURL, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse collection path: %w", err)
	}
	query := requestURL.Query()
	result := []T{}
	var previous json.RawMessage
	offset := 0
	for pageNumber := 0; pageNumber < listMaxPages; pageNumber++ {
		if paginated {
			query.Set("offset", strconv.Itoa(offset))
			query.Set("limit", strconv.Itoa(listPageSize))
			requestURL.RawQuery = query.Encode()
		}
		body, err := GetRaw(ctx, c, requestURL.String())
		if err != nil {
			return nil, err
		}
		page, err := decodeResponse[[]T](path, body)
		if err != nil {
			return nil, err
		}
		var envelope struct {
			Items      json.RawMessage `json:"items"`
			Offset     *int            `json:"offset"`
			TotalCount *int            `json:"totalCount"`
		}
		// Array responses have no pagination metadata.
		if bytes.HasPrefix(bytes.TrimSpace(body), []byte("{")) {
			if err := json.Unmarshal(body, &envelope); err != nil {
				return nil, fmt.Errorf("decode collection pagination: %w", err)
			}
		}
		if page == nil || len(*page) == 0 {
			if envelope.TotalCount != nil && offset < *envelope.TotalCount {
				return nil, fmt.Errorf("collection %q ended at offset %d before totalCount %d", path, offset, *envelope.TotalCount)
			}
			return &result, nil
		}
		if offset > 0 && bytes.Equal(previous, envelope.Items) && len(previous) > 0 {
			if envelope.TotalCount != nil || envelope.Offset != nil {
				return nil, fmt.Errorf("collection %q repeated a page at offset %d", path, offset)
			}
			// Older APIs may ignore offset/limit and return the entire collection.
			return &result, nil
		}
		if envelope.Offset != nil && *envelope.Offset != offset {
			return nil, fmt.Errorf("collection %q returned offset %d, requested %d", path, *envelope.Offset, offset)
		}
		result = append(result, (*page)...)
		offset += len(*page)
		if !paginated || bytes.HasPrefix(bytes.TrimSpace(body), []byte("[")) || (envelope.TotalCount != nil && offset >= *envelope.TotalCount) {
			return &result, nil
		}
		previous = envelope.Items
	}
	return nil, fmt.Errorf("collection %q exceeded the pagination limit of %d pages", path, listMaxPages)
}
