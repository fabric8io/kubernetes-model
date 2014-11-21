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
	Type string `json:"type"`
}

type JSONObjectDescriptor struct {
	Properties           map[string]JSONPropertyDescriptor `json:"properties,omitempty"`
	Required             []string                          `json:"required,omitempty"`
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

type JSONPropertyDescriptor struct {
	*JSONDescriptor
	*JSONReferenceDescriptor
	*JSONObjectDescriptor
	*JSONArrayDescriptor
	*JSONMapDescriptor
	*JavaTypeDescriptor
}

type JSONMapDescriptor struct {
	MapValueType JSONPropertyDescriptor `json:"additionalProperty"`
}
