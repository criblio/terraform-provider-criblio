package parser

// ResourceDef describes one Terraform resource discovered from OpenAPI annotations.
type ResourceDef struct {
	Name           string
	FileStem       string
	TypeName       string
	StructName     string
	SchemaName     string
	Create         OperationDef
	Read           OperationDef
	Update         OperationDef
	Delete         OperationDef
	List           OperationDef
	ListName       string
	ListFileStem   string
	ListStructName string
	ListTypeName   string
	Fields         []FieldDef
	OneOfVariants  []OneOfVariantDef
	Outputs        []OutputFile
	Action         bool
	NoRead         bool
}

// OperationDef describes one annotated OpenAPI operation.
type OperationDef struct {
	Method                   string
	Path                     string
	OperationID              string
	RequestSchema            string
	ResponseSchema           string
	PathParams               []FieldDef
	QueryParams              []FieldDef
	Examples                 []ExampleDef
	ReadAfterWrite           bool
	PreserveInputsAfterWrite bool
	ResetBody                any
	DeleteHook               string
}

// ExampleDef describes a request-body example attached to an OpenAPI operation.
type ExampleDef struct {
	Name    string
	Summary string
	Value   any
}

// FieldDef describes one Terraform model field.
type FieldDef struct {
	APIName            string
	TerraformName      string
	GoName             string
	Type               string
	ElementType        string
	Description        string
	NestedModelName    string
	NestedAPIModelName string
	NestedAttrTypes    string
	Required           bool
	Optional           bool
	Computed           bool
	OptionalComputed   bool
	Sensitive          bool
	PreferState        bool
	SuppressDiff       bool
	ForceNew           bool
	Ignored            bool
	CustomType         string
	ElementCustomType  string
	ReadOnly           bool
	WriteOnly          bool
	PathParam          bool
	QueryParam         bool
	RequestField       bool
	UpdateField        bool
	ApplyStrategy      string
	PlanModifierHook   string
	UseStateForUnknown bool
	EmitEmpty          bool
	FixedValue         string
	Enum               []string
	Fields             []FieldDef
	ObjectAsJSON       bool
	NotNull            bool
	ValidJSON          bool
	PipelineFunctionID bool
}

// OneOfVariantDef describes one flattened oneOf variant model.
type OneOfVariantDef struct {
	APIName            string
	TerraformName      string
	GoName             string
	ModelName          string
	SchemaName         string
	DiscriminatorValue string
	Fields             []FieldDef
}

// OutputFile describes a generated file decision.
type OutputFile struct {
	Path    string
	Kind    string
	Skipped bool
}
