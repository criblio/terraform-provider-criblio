// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type ExprLibEntry struct {
	Context  map[string]any `json:"context,omitempty"`
	EvalType *string        `json:"evalType,omitempty"`
	// JavaScript expression to evaluate
	Expr        string         `json:"expr"`
	ID          string         `json:"id"`
	Pack        *string        `json:"pack,omitempty"`
	Result      map[string]any `json:"result,omitempty"`
	Unprotected *bool          `json:"unprotected,omitempty"`
}

func (o *ExprLibEntry) GetContext() map[string]any {
	if o == nil {
		return nil
	}
	return o.Context
}

func (o *ExprLibEntry) GetEvalType() *string {
	if o == nil {
		return nil
	}
	return o.EvalType
}

func (o *ExprLibEntry) GetExpr() string {
	if o == nil {
		return ""
	}
	return o.Expr
}

func (o *ExprLibEntry) GetID() string {
	if o == nil {
		return ""
	}
	return o.ID
}

func (o *ExprLibEntry) GetPack() *string {
	if o == nil {
		return nil
	}
	return o.Pack
}

func (o *ExprLibEntry) GetResult() map[string]any {
	if o == nil {
		return nil
	}
	return o.Result
}

func (o *ExprLibEntry) GetUnprotected() *bool {
	if o == nil {
		return nil
	}
	return o.Unprotected
}
