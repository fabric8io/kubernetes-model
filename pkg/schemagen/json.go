package schemagen

type JSONSchema struct {
	ID          string                            `json:"id"`
	Schema      string                            `json:"$schema"`
	Description string                            `json:"description,omitempty"`
	Definitions map[string]JSONPropertyDescriptor `json:"definitions"`
	JSONDescriptor
	*JSONObjectDescriptor
}

type JSONDescriptor struct {
	Type        string        `json:"type"`
	Description string        `json:"description"`
	Default     string        `json:"default,omitempty"`
	Required    bool          `json:"required,omitempty"`
	Minimum     float32       `json:"minimum,omitempty"`
	Maximum     float32       `json:"maximum,omitempty"`
	MinItems    int           `json:"minItems,omitempty"`
	MaxItems    int           `json:"maxItems,omitempty"`
	MinLength   int           `json:"minLength,omitempty"`
	MaxLength   int           `json:"maxLength,omitempty"`
	Pattern     string        `json:"pattern,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
}

type JSONObjectDescriptor struct {
	Properties           map[string]JSONPropertyDescriptor `json:"properties,omitempty"`
	AdditionalProperties bool                              `json:"additionalProperties"`
}

type JSONArrayDescriptor struct {
	Items JSONPropertyDescriptor `json:"items"`
}

type JSONReferenceDescriptor struct {
	Reference string `json:"$ref"`
}

type JavaTypeDescriptor struct {
	JavaType string `json:"javaType"`
}

type JavaInterfacesDescriptor struct {
	JavaInterfaces []string `json:"javaInterfaces,omitempty"`
}

type JSONPropertyDescriptor struct {
	*JSONDescriptor
	*JSONReferenceDescriptor
	*JSONObjectDescriptor
	*JSONArrayDescriptor
	*JSONMapDescriptor
	*JavaTypeDescriptor
	*JavaInterfacesDescriptor
}

type JSONMapDescriptor struct {
	MapValueType JSONPropertyDescriptor `json:"additionalProperty"`
}
