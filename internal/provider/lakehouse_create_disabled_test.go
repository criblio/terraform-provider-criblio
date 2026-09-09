package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestCriblLakeHouseCreateDisabled(t *testing.T) {
	var resp resource.CreateResponse

	(&CriblLakeHouseResource{}).Create(context.Background(), resource.CreateRequest{}, &resp)

	assertCreateDisabledDiagnostic(t, resp)
}

func assertCreateDisabledDiagnostic(t *testing.T, resp resource.CreateResponse) {
	t.Helper()
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected resource creation to be rejected")
	}
	detail := resp.Diagnostics[0].Detail()
	if !strings.Contains(detail, "Existing") || !strings.Contains(detail, "Cribl Search") {
		t.Fatalf("expected migration guidance in diagnostic, got %q", detail)
	}
}
