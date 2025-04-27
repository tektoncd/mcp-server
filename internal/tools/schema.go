package tools

import (
	"encoding/json"

	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	common "k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/validation/spec"
)

func getSchemaForType(name string) string {
	return getSchemaForTypeWithDefinitions(name, v1.GetOpenAPIDefinitions(
		func(path string) spec.Ref {
			return spec.MustCreateRef(path)
		},
	))
}

func getSchemaForTypeWithDefinitions(name string, definitions map[string]common.OpenAPIDefinition) string {
	definition, ok := definitions[name]
	if !ok {
		return ""
	}
	resolvedDefinition := resolveReferences(definition.Schema, definitions)
	data, _ := json.Marshal(resolvedDefinition)
	return string(data)
}

func resolveReferences(schema spec.Schema, definitions map[string]common.OpenAPIDefinition) spec.Schema {
	if schema.Type.Contains("object") {
		for name, prop := range schema.Properties {
			schema.Properties[name] = resolveReferences(prop, definitions)
		}
	}
	if schema.Type.Contains("array") {
		if schema.Items.Len() == 1 {
			*schema.Items.Schema = resolveReferences(*schema.Items.Schema, definitions)
		} else {
			resolvedSchemas := make([]spec.Schema, 0)
			for _, sch := range schema.Items.Schemas {
				resolvedSchemas = append(resolvedSchemas, resolveReferences(sch, definitions))
			}
			schema.Items.Schemas = resolvedSchemas
		}
	}
	if path := schema.Ref.String(); path != "" {
		definition := definitions[path]
		return resolveReferences(definition.Schema, definitions)
	}
	return schema
}
