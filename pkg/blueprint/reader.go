package blueprint

import (
	"github.com/cruciblehq/protocol/pkg/codec"
)

// Loads a blueprint from a file.
//
// The path parameter specifies the full path to the blueprint file. The file
// format is inferred from the extension (.yaml, .json, .toml).
func Read(path string) (*Blueprint, error) {
	var bp Blueprint
	if _, err := codec.DecodeFile(path, "field", &bp); err != nil {
		return nil, err
	}
	return &bp, nil
}
