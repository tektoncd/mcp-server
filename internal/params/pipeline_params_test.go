package params

import (
	"reflect"
	"testing"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

func TestParsePipelineRunParams(t *testing.T) {
	tests := []struct {
		name      string
		paramStr  string
		want      []pipelinev1.Param
		wantError bool
	}{
		{
			name:     "empty string",
			paramStr: "",
			want:     nil,
		},
		{
			name:     "single string param",
			paramStr: "message=Hello",
			want: []pipelinev1.Param{
				{
					Name: "message",
					Value: pipelinev1.ParamValue{
						Type:      pipelinev1.ParamTypeString,
						StringVal: "Hello",
					},
				},
			},
		},
		{
			name:     "multiple string params",
			paramStr: "message=Hello,user=World",
			want: []pipelinev1.Param{
				{
					Name: "message",
					Value: pipelinev1.ParamValue{
						Type:      pipelinev1.ParamTypeString,
						StringVal: "Hello",
					},
				},
				{
					Name: "user",
					Value: pipelinev1.ParamValue{
						Type:      pipelinev1.ParamTypeString,
						StringVal: "World",
					},
				},
			},
		},
		{
			name:     "array param",
			paramStr: "items=array:item1:item2:item3",
			want: []pipelinev1.Param{
				{
					Name: "items",
					Value: pipelinev1.ParamValue{
						Type:     pipelinev1.ParamTypeArray,
						ArrayVal: []string{"item1", "item2", "item3"},
					},
				},
			},
		},
		{
			name:     "object param",
			paramStr: "config=object:key1=value1:key2=value2",
			want: []pipelinev1.Param{
				{
					Name: "config",
					Value: pipelinev1.ParamValue{
						Type: pipelinev1.ParamTypeObject,
						ObjectVal: map[string]string{
							"key1": "value1",
							"key2": "value2",
						},
					},
				},
			},
		},
		{
			name:     "mixed param types",
			paramStr: "message=Hello,items=array:item1:item2,config=object:key1=value1:key2=value2",
			want: []pipelinev1.Param{
				{
					Name: "message",
					Value: pipelinev1.ParamValue{
						Type:      pipelinev1.ParamTypeString,
						StringVal: "Hello",
					},
				},
				{
					Name: "items",
					Value: pipelinev1.ParamValue{
						Type:     pipelinev1.ParamTypeArray,
						ArrayVal: []string{"item1", "item2"},
					},
				},
				{
					Name: "config",
					Value: pipelinev1.ParamValue{
						Type: pipelinev1.ParamTypeObject,
						ObjectVal: map[string]string{
							"key1": "value1",
							"key2": "value2",
						},
					},
				},
			},
		},
		{
			name:      "invalid format - missing value",
			paramStr:  "message",
			wantError: true,
		},
		{
			name:     "invalid object format",
			paramStr: "config=object:invalid",
			want: []pipelinev1.Param{
				{
					Name: "config",
					Value: pipelinev1.ParamValue{
						Type:      pipelinev1.ParamTypeString,
						StringVal: "object:invalid",
					},
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePipelineRunParams(tt.paramStr)
			if (err != nil) != tt.wantError {
				t.Errorf("ParsePipelineRunParams() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParsePipelineRunParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMergeParams(t *testing.T) {
	existing := []pipelinev1.Param{
		{
			Name: "message",
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: "Hello",
			},
		},
		{
			Name: "user",
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: "World",
			},
		},
	}

	new := []pipelinev1.Param{
		{
			Name: "message",
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: "Updated",
			},
		},
		{
			Name: "items",
			Value: pipelinev1.ParamValue{
				Type:     pipelinev1.ParamTypeArray,
				ArrayVal: []string{"item1", "item2"},
			},
		},
	}

	expected := []pipelinev1.Param{
		{
			Name: "message",
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: "Updated",
			},
		},
		{
			Name: "user",
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: "World",
			},
		},
		{
			Name: "items",
			Value: pipelinev1.ParamValue{
				Type:     pipelinev1.ParamTypeArray,
				ArrayVal: []string{"item1", "item2"},
			},
		},
	}

	result := MergeParams(existing, new)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("MergeParams() = %v, want %v", result, expected)
	}

	// Test empty params
	if !reflect.DeepEqual(MergeParams(nil, new), new) {
		t.Errorf("MergeParams(nil, new) should return new")
	}

	if !reflect.DeepEqual(MergeParams(existing, nil), existing) {
		t.Errorf("MergeParams(existing, nil) should return existing")
	}
}
