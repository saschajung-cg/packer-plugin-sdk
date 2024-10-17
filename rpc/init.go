// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

/*
Package rpc contains the implementation of the remote procedure call code that
the Packer core uses to communicate with packer plugins. As a plugin maintainer,
you are unlikely to need to directly import or use this package, but it
underpins the packer server that all plugins must implement.
*/
package rpc

import (
	"bytes"
	"encoding/gob"

	"github.com/zclconf/go-cty/cty"

	"encoding/json"
)

// Test that cty types implement the gob.GobEncoder interface.
// Support for encoding/gob was removed in github.com/zclconf/go-cty@v1.11.0.
// Refer to issue https://github.com/hashicorp/packer-plugin-sdk/issues/187

func (v cty.Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.v) // Serialize the internal value v
}

func (v *cty.Value) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.v) // Deserialize back to the internal value
}

func init() {
	gob.Register(new(map[string]string))
	gob.Register(make([]interface{}, 0))
	gob.Register(new(BasicError))
}

// Implement the GobEncode method for cty.Value
func (v cty.Value) GobEncode() ([]byte, error) {
	// Manually serialize the internal value v using gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v.v) // Encode the internal value
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Implement the GobDecode method for cty.Value
func (v *cty.Value) GobDecode(data []byte) error {
	// Manually decode the internal value v using gob
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(&v.v) // Decode back to the internal value
}
