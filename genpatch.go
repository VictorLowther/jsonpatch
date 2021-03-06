package jsonpatch

import (
	"encoding/json"
	"reflect"

	"github.com/VictorLowther/jsonpatch/utils"
)

// This generator does not create copy or move patch ops, and I don't
// care enough to optimize it to do so.  Ditto for slice handling.
// There is a lot of optimization that could be done here, but it can get complex real quick.
func basicGen(base, target interface{}, paranoid bool, ptr pointer) patch {
	res := make(patch, 0)
	if reflect.TypeOf(base) != reflect.TypeOf(target) {
		if paranoid {
			res = append(res, operation{"test", ptr, nil, utils.Clone(base)})
		}
		res = append(res, operation{"replace", ptr, nil, utils.Clone(target)})
		return res
	}
	switch baseVal := base.(type) {
	case map[string]interface{}:
		targetVal := target.(map[string]interface{})
		handled := make(map[string]struct{})
		// Handle removed and changed first.
		for k, oldVal := range baseVal {
			newPtr := ptr.Append(k)
			newVal, ok := targetVal[k]
			if !ok {
				// Generate a remove op
				if paranoid {
					res = append(res, operation{"test", newPtr, nil, utils.Clone(oldVal)})
				}
				res = append(res, operation{"remove", newPtr, nil, nil})
			} else {
				subPatch := basicGen(oldVal, newVal, paranoid, newPtr)
				res = append(res, subPatch...)
			}
			handled[k] = struct{}{}
		}
		// Now, handle additions
		for k, newVal := range targetVal {
			if _, ok := handled[k]; ok {
				continue
			}
			res = append(res, operation{"add", ptr.Append(k), nil, utils.Clone(newVal)})
		}
	// case []interface{}:
	// Eventually, add code to handle slices more
	// efficiently.  For now, through, be dumb.
	default:
		if !reflect.DeepEqual(base, target) {
			if paranoid {
				res = append(res, operation{"test", ptr, nil, utils.Clone(base)})
			}
			res = append(res, operation{"replace", ptr, nil, utils.Clone(target)})
		}
	}
	return res
}

// Generate generates a JSON Patch that will modify base into target.
// If paranoid is true, then the generated patch will have test checks.
//
// base and target must be the result of unmarshalling JSON into an interface{}
func Generate(base, target interface{}, paranoid bool) ([]byte, error) {
	p := basicGen(base, target, paranoid, make(pointer, 0))
	return json.Marshal(p)
}

// GenerateJSON does the same thing as Generate, except base and
// target should be byte arrays containing raw JSON
func GenerateJSON(base, target []byte, paranoid bool) ([]byte, error) {
	var rawBase, rawTarget interface{}
	if err := json.Unmarshal(base, &rawBase); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(target, &rawTarget); err != nil {
		return nil, err
	}
	return Generate(rawBase, rawTarget, paranoid)
}
