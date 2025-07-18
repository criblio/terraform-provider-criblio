// Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.

package shared

type PackClone struct {
	Dest       *string  `json:"dest,omitempty"`
	DstGroups  []string `json:"dstGroups"`
	Force      *bool    `json:"force,omitempty"`
	IsDisabled *bool    `json:"isDisabled,omitempty"`
	Packs      []string `json:"packs"`
	SrcGroup   string   `json:"srcGroup"`
}

func (o *PackClone) GetDest() *string {
	if o == nil {
		return nil
	}
	return o.Dest
}

func (o *PackClone) GetDstGroups() []string {
	if o == nil {
		return []string{}
	}
	return o.DstGroups
}

func (o *PackClone) GetForce() *bool {
	if o == nil {
		return nil
	}
	return o.Force
}

func (o *PackClone) GetIsDisabled() *bool {
	if o == nil {
		return nil
	}
	return o.IsDisabled
}

func (o *PackClone) GetPacks() []string {
	if o == nil {
		return []string{}
	}
	return o.Packs
}

func (o *PackClone) GetSrcGroup() string {
	if o == nil {
		return ""
	}
	return o.SrcGroup
}
