package utils

import "github.com/manicminer/hamilton/msgraph"

func Bool(input bool) *bool {
	return &input
}

func Int32(input int32) *int32 {
	return &input
}

func String(input string) *string {
	return &input
}

func NullableString(input string) *msgraph.StringNullWhenEmpty {
	output := msgraph.StringNullWhenEmpty(input)
	return &output
}
