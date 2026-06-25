package custom

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseCriblLakeDatasetListBodySkipsDeletedDatasets(t *testing.T) {
	body := []byte(`{
		"items": [
			{"id": "ready_dataset", "status": "ready"},
			{"id": "deleting_dataset", "status": "marked_for_deletion"},
			{"id": "terminated_dataset", "status": "terminated"},
			{"id": "flagged_dataset", "markedForDeletion": true},
			{"id": "timestamped_dataset", "deletionStartedAt": 1781874901035},
			{"id": "cribl_logs", "status": "ready"}
		]
	}`)

	ids, err := ParseCriblLakeDatasetListBody(body, "default")
	require.NoError(t, err)

	require.Len(t, ids, 1)
	assert.Equal(t, "ready_dataset", ids[0]["id"])
	assert.Equal(t, "default", ids[0]["lake_id"])
}

func TestCriblLakeDatasetListKey(t *testing.T) {
	assert.Equal(t, "cribl_lake_dataset_list:default", criblLakeDatasetListKey("/api/v1/products/lake/lakes/default/datasets"))
	assert.Empty(t, criblLakeDatasetListKey("/api/v1/products/lake/lakes/default/datasets/my_dataset"))
}
