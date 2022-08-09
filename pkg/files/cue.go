package files

import (
	"encoding/json"
	"os"

	"cuelang.org/go/cue/cuecontext"
)

func ParseCUEFile(source string, obj interface{}) error {
	bytesRead, err := os.ReadFile(source)

	if err != nil {
		return err
	}
	ctx := cuecontext.New()
	jsonBytes, err := ctx.CompileString(string(bytesRead)).MarshalJSON()
	if err != nil {
		return err
	}
	// This unmarshalling cannot give error because if there was a problem
	// with the cue parsing it would have bombed in the CompileString
	_ = json.Unmarshal(jsonBytes, obj)

	return nil
}

// // UnmarshalTo unmarshal value into golang object
// func (val cue.Value) UnmarshalTo(x interface{}) error {
// 	data, err := val.MarshalJSON()
// 	if err != nil {
// 		return err
// 	}
// 	return json.Unmarshal(data, x)
// }
