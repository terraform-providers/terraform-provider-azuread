// Code generated by "stringer -type=ResourceMode -output=resource_mode_string.go resource_mode.go"; DO NOT EDIT.

package config

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ManagedResourceMode-0]
	_ = x[DataResourceMode-1]
}

const _ResourceMode_name = "ManagedResourceModeDataResourceMode"

var _ResourceMode_index = [...]uint8{0, 19, 35}

func (i ResourceMode) String() string {
	if i < 0 || i >= ResourceMode(len(_ResourceMode_index)-1) {
		return "ResourceMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ResourceMode_name[_ResourceMode_index[i]:_ResourceMode_index[i+1]]
}
