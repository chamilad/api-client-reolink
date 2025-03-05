// Package models contains the models that can be used to communciate with the camera API
package models

import "encoding/json"

// yes, I know utils is bad, but it's named common, so it's all good
// todo: find a better place for these

// BoolFromInt is a type to unmarshal ints as bools
// camera API has been observed to send 1 as true and 0 as false
type BoolFromInt bool

// UnmarshalJSON is the custom unmarshaller to convert ints to bools
func (b *BoolFromInt) UnmarshalJSON(data []byte) error {
	var temp int
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	*b = BoolFromInt(temp == 1)
	return nil
}
