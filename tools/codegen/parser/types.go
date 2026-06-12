package parser

// ResourceDef describes one Terraform resource discovered from OpenAPI annotations.
type ResourceDef struct {
	Name          string
	TypeName      string
	StructName    string
	SchemaName    string
	Create        OperationDef
	Read          OperationDef
	Update        OperationDef
	Delete        OperationDef
	Fields        []FieldDef
	OneOfVariants []OneOfVariantDef
	Outputs       []OutputFile
}

// OperationDef describes one annotated OpenAPI operation.
type OperationDef struct {
	Method         string
	Path           string
	OperationID    string
	RequestSchema  string
	ResponseSchema string
	PathParams     []FieldDef
}

// FieldDef describes one Terraform model field.
type FieldDef struct {
	APIName       string
	TerraformName string
	GoName        string
	Type          string
	ElementType   string
	Description   string
	Required      bool
	Optional      bool
	Computed      bool
	Sensitive     bool
	PreferState   bool
	ForceNew      bool
	Ignored       bool
	CustomType    string
	ReadOnly      bool
	WriteOnly     bool
	PathParam     bool
	RequestField  bool
	ApplyStrategy string
	Fields        []FieldDef
}

// OneOfVariantDef describes one flattened oneOf variant model.
type OneOfVariantDef struct {
	APIName       string
	TerraformName string
	GoName        string
	ModelName     string
	SchemaName    string
	Fields        []FieldDef
}

// OutputFile describes a generated file decision.
type OutputFile struct {
	Path    string
	Kind    string
	Skipped bool
}
