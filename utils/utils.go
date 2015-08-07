package utils

// Holds a couple of useful utilities for JSON handling

import (
	"encoding/json"
	"reflect"
)

// Clone performs a deep clone of a JSON-ish structure.
func Clone(val interface{}) interface{} {
	switch t := val.(type) {
	case []interface{}:
		res := make([]interface{}, len(t))
		for i := range t {
			res[i] = Clone(t[i])
		}
		return res
	case map[string]interface{}:
		res := make(map[string]interface{}, len(t))
		for k, v := range t {
			res[k] = Clone(v)
		}
		return res
	default:
		return val
	}
}

func merge(src, changes interface{}) interface{} {
	if reflect.TypeOf(src) != reflect.TypeOf(changes) {
		return changes
	}
	switch srcVal := src.(type) {
	case map[string]interface{}:
		changesVal := changes.(map[string]interface{})
		for k, newVal := range changesVal {
			if newVal == nil {
				delete(srcVal, k)
				continue
			} else {
				srcVal[k] = Merge(srcVal[k], newVal)
			}
		}
		return srcVal
	default:
		return changes
	}
}

// Merge merges changes into src recursively.  The original objects
// will be left unchanged.
func Merge(src, changes interface{}) interface{} {
	return merge(Clone(src), Clone(changes))
}

// MergeJSON does the same as Merge, except it accepts and returns
// byte arrays that
func MergeJSON(src, changes []byte) ([]byte, error) {
	var srcObj, changesObj, resObj interface{}
	if err := json.Unmarshal(src, &srcObj); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(changes, &changesObj); err != nil {
		return nil, err
	}
	resObj = merge(srcObj, changesObj)
	return json.Marshal(resObj)
}

// Remarshal marshals src and then unmarshals it into target.
func Remarshal(src, target interface{}) error {
	r, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(r, &target)
}
