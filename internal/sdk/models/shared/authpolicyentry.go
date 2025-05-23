// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type AuthPolicyEntry struct {
	Actions []string `json:"actions"`
	Object  string   `json:"object"`
}

func (o *AuthPolicyEntry) GetActions() []string {
	if o == nil {
		return []string{}
	}
	return o.Actions
}

func (o *AuthPolicyEntry) GetObject() string {
	if o == nil {
		return ""
	}
	return o.Object
}
