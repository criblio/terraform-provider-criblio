package converter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiagnosticsToCLIErrors(t *testing.T) {
	t.Run("nil_empty_returns_nil", func(t *testing.T) {
		assert.Nil(t, DiagnosticsToCLIErrors(nil, "", ""))
		assert.Nil(t, DiagnosticsToCLIErrors(diag.Diagnostics{}, "criblio_source", "id1"))
	})

	t.Run("single_error_includes_summary_and_context", func(t *testing.T) {
		var diags diag.Diagnostics
		diags.AddError("Bad value", "field X is invalid")
		out := DiagnosticsToCLIErrors(diags, "criblio_source", "default/my-id")
		require.Len(t, out, 1)
		assert.Equal(t, "criblio_source", out[0].ResourceType)
		assert.Equal(t, "default/my-id", out[0].ResourceID)
		assert.Equal(t, "Bad value", out[0].Summary)
		assert.Equal(t, "field X is invalid", out[0].Detail)
	})

	t.Run("multiple_errors_all_converted", func(t *testing.T) {
		var diags diag.Diagnostics
		diags.AddError("First", "detail one")
		diags.AddError("Second", "detail two")
		out := DiagnosticsToCLIErrors(diags, "criblio_pipeline", "")
		require.Len(t, out, 2)
		assert.Equal(t, "First", out[0].Summary)
		assert.Equal(t, "Second", out[1].Summary)
	})

	t.Run("warning_included", func(t *testing.T) {
		var diags diag.Diagnostics
		diags.AddWarning("Warning", "something to note")
		out := DiagnosticsToCLIErrors(diags, "", "")
		require.Len(t, out, 1)
		assert.Equal(t, "Warning", out[0].Summary)
	})
}

func TestDiagnosticsToError(t *testing.T) {
	t.Run("no_errors_returns_nil", func(t *testing.T) {
		assert.NoError(t, DiagnosticsToError(nil, "", ""))
		var diags diag.Diagnostics
		diags.AddWarning("Only warning", "")
		assert.NoError(t, DiagnosticsToError(diags, "", ""))
	})

	t.Run("single_error_returns_combined_message", func(t *testing.T) {
		var diags diag.Diagnostics
		diags.AddError("Summary", "Detail")
		err := DiagnosticsToError(diags, "criblio_source", "id1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "criblio_source")
		assert.Contains(t, err.Error(), "id1")
		assert.Contains(t, err.Error(), "Summary")
	})

	t.Run("multiple_errors_joined", func(t *testing.T) {
		var diags diag.Diagnostics
		diags.AddError("A", "a")
		diags.AddError("B", "b")
		err := DiagnosticsToError(diags, "t", "id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "A")
		assert.Contains(t, err.Error(), "B")
	})
}

func TestCLIError_Error(t *testing.T) {
	t.Run("with_resource_context", func(t *testing.T) {
		e := CLIError{ResourceType: "criblio_source", ResourceID: "g/id", Summary: "failed"}
		assert.Contains(t, e.Error(), "criblio_source")
		assert.Contains(t, e.Error(), "g/id")
		assert.Contains(t, e.Error(), "failed")
	})
	t.Run("summary_only", func(t *testing.T) {
		e := CLIError{Summary: "something broke"}
		assert.Equal(t, "something broke", e.Error())
	})
}
