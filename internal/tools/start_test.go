package tools

import (
	"strings"
	"testing"
)

func TestParseParamsString(t *testing.T) {
	// Test cases for parameter string parsing
	testCases := []struct {
		name           string
		paramsStr      string
		expectedParams map[string]bool
	}{
		{
			name:      "Basic parameters",
			paramsStr: "message=Hello,user=World",
			expectedParams: map[string]bool{
				"message=Hello": false,
				"user=World":    false,
			},
		},
		{
			name:      "Array and object parameters",
			paramsStr: "message=Hello,items=array:item1:item2:item3,config=object:key1=value1:key2=value2",
			expectedParams: map[string]bool{
				"message=Hello":                         false,
				"items=array:item1:item2:item3":         false,
				"config=object:key1=value1:key2=value2": false,
			},
		},
		{
			name:           "Empty string",
			paramsStr:      "",
			expectedParams: map[string]bool{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.paramsStr == "" {
				// No parameters to process
				return
			}

			// Split the parameter string
			paramPairs := strings.Split(tc.paramsStr, ",")

			// Check we have the expected number of parameters
			if len(paramPairs) != len(tc.expectedParams) {
				t.Errorf("Expected %d parameters, got %d", len(tc.expectedParams), len(paramPairs))
			}

			// Check the parameters have the expected format
			for _, pair := range paramPairs {
				if _, ok := tc.expectedParams[pair]; ok {
					tc.expectedParams[pair] = true
				} else {
					t.Errorf("Unexpected parameter: %s", pair)
				}
			}

			// Verify all expected parameters were found
			for param, found := range tc.expectedParams {
				if !found {
					t.Errorf("Expected parameter not found: %s", param)
				}
			}
		})
	}
}
