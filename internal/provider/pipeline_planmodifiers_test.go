package provider

import "testing"

func TestPipelineFunctionConfEquivalence(t *testing.T) {
	for _, tc := range []struct {
		name, id, config, api string
		equal                 bool
	}{
		{"eval default", "eval", `{"add":[]}`, `{"printUndefineds":false,"add":[]}`, true},
		{"aggregation default", "aggregation", `{}`, `{"printUndefineds":false}`, true},
		{"enabled default", "eval", `{}`, `{"printUndefineds":true}`, false},
		{"other function", "redis", `{}`, `{"printUndefineds":false}`, false},
		{"nested user key", "eval", `{"obj":{}}`, `{"obj":{"printUndefineds":false}}`, false},
		{"key order", "redis", `{"command":"get","key":"host"}`, `{"key":"host","command":"get"}`, true},
		{"command change", "redis", `{"command":"get"}`, `{"command":"set"}`, false},
		{"array order", "eval", `{"remove":["a","b"]}`, `{"remove":["b","a"]}`, false},
		{"numeric formatting", "eval", `{"n":1.0}`, `{"n":1}`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			equal := normalizePipelineConf(tc.id, tc.config) == normalizePipelineConf(tc.id, tc.api)
			if equal != tc.equal {
				t.Errorf("config/API equality = %t, want %t", equal, tc.equal)
			}
		})
	}
}
