package params

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
)

// ExtractPipelineRunParams extracts parameters from the request arguments and returns a slice of Param objects
func ExtractPipelineRunParams(r mcp.ReadResourceRequest) ([]pipelinev1.Param, error) {
	paramsArg, ok := r.Params.Arguments["params"].([]string)
	if !ok || len(paramsArg) == 0 {
		return nil, nil
	}

	paramStr := paramsArg[0]
	if paramStr == "" {
		return nil, nil
	}

	return ParsePipelineRunParams(paramStr)
}

// ParsePipelineRunParams parses a string representation of parameters into a slice of Param objects
// Format: key1=value1,key2=value2,key3=array:val1:val2:val3,key4=object:k1=v1:k2=v2
func ParsePipelineRunParams(paramStr string) ([]pipelinev1.Param, error) {
	if paramStr == "" {
		return nil, nil
	}

	paramPairs := strings.Split(paramStr, ",")
	params := make([]pipelinev1.Param, 0, len(paramPairs))

	for _, pair := range paramPairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid parameter format: %s", pair)
		}

		key := kv[0]
		value := kv[1]

		param, err := parseParamValue(key, value)
		if err != nil {
			slog.Warn(fmt.Sprintf("Error parsing parameter %s: %v", key, err))
			// Create a string parameter with the original value
			param = pipelinev1.Param{
				Name: key,
				Value: pipelinev1.ParamValue{
					Type:      pipelinev1.ParamTypeString,
					StringVal: value,
				},
			}
		}

		params = append(params, param)
	}

	return params, nil
}

// parseParamValue converts a string value to the appropriate ParamValue type
func parseParamValue(key, value string) (pipelinev1.Param, error) {
	if strings.HasPrefix(value, "array:") {
		// Handle array type param
		arrayValues := strings.Split(strings.TrimPrefix(value, "array:"), ":")
		return pipelinev1.Param{
			Name: key,
			Value: pipelinev1.ParamValue{
				Type:     pipelinev1.ParamTypeArray,
				ArrayVal: arrayValues,
			},
		}, nil
	} else if strings.HasPrefix(value, "object:") {
		// Handle object type param
		objectStr := strings.TrimPrefix(value, "object:")
		objectPairs := strings.Split(objectStr, ":")
		objectVal := make(map[string]string)

		// If we don't have any valid key-value pairs, return an error
		validPairs := false

		for _, objectPair := range objectPairs {
			objKV := strings.SplitN(objectPair, "=", 2)
			if len(objKV) != 2 {
				continue
			}
			objectVal[objKV[0]] = objKV[1]
			validPairs = true
		}

		if !validPairs {
			return pipelinev1.Param{}, fmt.Errorf("invalid object format, no valid key-value pairs found")
		}

		return pipelinev1.Param{
			Name: key,
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeObject,
				ObjectVal: objectVal,
			},
		}, nil
	} else {
		// Handle string type param
		return pipelinev1.Param{
			Name: key,
			Value: pipelinev1.ParamValue{
				Type:      pipelinev1.ParamTypeString,
				StringVal: value,
			},
		}, nil
	}
}

// MergeParams merges new parameters with existing ones, either replacing or adding them
func MergeParams(existing []pipelinev1.Param, new []pipelinev1.Param) []pipelinev1.Param {
	if len(existing) == 0 {
		return new
	}

	if len(new) == 0 {
		return existing
	}

	// Create a map of existing params for easier lookup/replacement
	existingParams := make(map[string]int)
	for i, p := range existing {
		existingParams[p.Name] = i
	}

	result := make([]pipelinev1.Param, len(existing))
	copy(result, existing)

	// Update or append parameters
	for _, newParam := range new {
		if idx, exists := existingParams[newParam.Name]; exists {
			// Replace existing parameter
			result[idx] = newParam
		} else {
			// Append new parameter
			result = append(result, newParam)
		}
	}

	return result
}
