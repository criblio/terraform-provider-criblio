package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSearchDatasetIDFallsBackToActiveVariant(t *testing.T) {
	model := SearchDatasetModel{
		ID: types.StringNull(),
		DatasetS3: &DatasetS3Model{
			ID: types.StringValue("s3-dataset-001"),
		},
	}

	if got := searchDatasetID(model); got != "s3-dataset-001" {
		t.Fatalf("searchDatasetID() = %q, want %q", got, "s3-dataset-001")
	}
}
