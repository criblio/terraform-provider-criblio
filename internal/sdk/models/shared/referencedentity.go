// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type ReferencedEntity struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (o *ReferencedEntity) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *ReferencedEntity) GetType() string {
	if o == nil {
		return ""
	}
	return o.Type
}
