package provider

import (
	"testing"

	tfTypes "github.com/criblio/terraform-provider-criblio/internal/provider/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestPromoteFirstInputUnionItemToTopLevel_Source(t *testing.T) {
	m := &SourceResourceModel{
		GroupID: types.StringValue("default"),
		ID:      types.StringValue("open_telemetry"),
		Items: []tfTypes.InputUnion1{
			{
				InputOpenTelemetry: &tfTypes.InputOpenTelemetry{
					Type: types.StringValue("open_telemetry"),
					Host: types.StringValue("0.0.0.0"),
					Port: types.Float64Value(4318),
				},
			},
		},
	}
	require.Nil(t, m.InputOpenTelemetry)
	PromoteFirstInputUnionItemToTopLevel(m)
	require.NotNil(t, m.InputOpenTelemetry)
	require.Equal(t, "open_telemetry", m.InputOpenTelemetry.Type.ValueString())
	require.Equal(t, float64(4318), m.InputOpenTelemetry.Port.ValueFloat64())
}

func TestPromoteFirstInputUnionItemToTopLevel_NoItems(t *testing.T) {
	m := &SourceResourceModel{ID: types.StringValue("x")}
	PromoteFirstInputUnionItemToTopLevel(m)
	require.Nil(t, m.InputOpenTelemetry)
}

func TestPromoteFirstInputUnionItemToTopLevel_EmptySlicesBecomeNil(t *testing.T) {
	m := &SourceResourceModel{
		Items: []tfTypes.InputUnion1{
			{
				InputKubeMetrics: &tfTypes.InputKubeMetrics{
					Type:       types.StringValue("kube_metrics"),
					ID:         types.StringValue("in_kube_metrics"),
					Metadata:   []tfTypes.ItemsTypeMetadata{},
					Streamtags: []types.String{},
					Connections: []tfTypes.ItemsTypeConnectionsOptional{},
				},
			},
		},
	}
	PromoteFirstInputUnionItemToTopLevel(m)
	require.NotNil(t, m.InputKubeMetrics)
	require.Nil(t, m.InputKubeMetrics.Metadata)
	require.Nil(t, m.InputKubeMetrics.Streamtags)
	require.Nil(t, m.InputKubeMetrics.Connections)
}

func TestPromoteFirstInputUnionItemToTopLevel_NestedAuthTokenEmptySlicesBecomeNil(t *testing.T) {
	m := &SourceResourceModel{
		Items: []tfTypes.InputUnion1{
			{
				InputSplunkHec: &tfTypes.InputSplunkHec{
					Type: types.StringValue("splunk_hec"),
					ID:   types.StringValue("in_splunk_hec"),
					AuthTokens: []tfTypes.InputSplunkHecAuthToken{
						{
							Token:                 types.StringValue("t"),
							AllowedIndexesAtToken: []types.String{},
							Metadata:              []tfTypes.ItemsTypeMetadata{},
						},
					},
				},
			},
		},
	}
	PromoteFirstInputUnionItemToTopLevel(m)
	require.NotNil(t, m.InputSplunkHec)
	require.Len(t, m.InputSplunkHec.AuthTokens, 1)
	require.Nil(t, m.InputSplunkHec.AuthTokens[0].AllowedIndexesAtToken)
	require.Nil(t, m.InputSplunkHec.AuthTokens[0].Metadata)
}
