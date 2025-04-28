package params

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// Test the generic Optional function which handles optional parameters
func TestOptional(t *testing.T) {
	testCases := []struct {
		name             string
		request          mcp.CallToolRequest
		param            string
		expectedValue    interface{}
		expectedErrorMsg string
	}{
		{
			name: "parameter exists and is of correct type (string)",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"test_param": "test_value",
					},
				},
			},
			param:         "test_param",
			expectedValue: "test_value",
		},
		{
			name: "parameter exists and is of correct type (int)",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"test_param": 42,
					},
				},
			},
			param:         "test_param",
			expectedValue: 42,
		},
		{
			name: "parameter exists and is of correct type (bool)",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"test_param": true,
					},
				},
			},
			param:         "test_param",
			expectedValue: true,
		},
		{
			name: "parameter does not exist",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{},
				},
			},
			param:         "test_param",
			expectedValue: "",
		},
		{
			name: "parameter exists but is of incorrect type",
			request: mcp.CallToolRequest{
				Params: struct {
					Name      string                 `json:"name"`
					Arguments map[string]interface{} `json:"arguments,omitempty"`
					Meta      *struct {
						ProgressToken mcp.ProgressToken `json:"progressToken,omitempty"`
					} `json:"_meta,omitempty"`
				}{
					Arguments: map[string]interface{}{
						"test_param": 123,
					},
				},
			},
			param:            "test_param",
			expectedErrorMsg: "parameter test_param is not of type string, is int",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test for string
			if tc.name == "parameter exists but is of incorrect type" {
				val, err := Optional[string](tc.request, tc.param)
				assert.Equal(t, "", val)
				if err != nil {
					assert.Contains(t, err.Error(), tc.expectedErrorMsg)
				} else {
					t.Errorf("Expected error but got nil")
				}
			} else if tc.name == "parameter exists and is of correct type (string)" {
				val, err := Optional[string](tc.request, tc.param)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, val)
			} else if tc.name == "parameter exists and is of correct type (int)" {
				val, err := Optional[int](tc.request, tc.param)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, val)
			} else if tc.name == "parameter exists and is of correct type (bool)" {
				val, err := Optional[bool](tc.request, tc.param)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedValue, val)
			} else if tc.name == "parameter does not exist" {
				val, err := Optional[string](tc.request, tc.param)
				assert.NoError(t, err)
				assert.Equal(t, "", val)
			}
		})
	}
}
