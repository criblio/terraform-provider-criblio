// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type Error struct {
	// Error message
	Message *string `json:"message,omitempty"`
}

func (o *Error) GetMessage() *string {
	if o == nil {
		return nil
	}
	return o.Message
}
