package tf

func ExpandStringSlice(input []interface{}) []string {
	result := make([]string, 0)
	for _, item := range input {
		result = append(result, item.(string))
	}
	return result
}

func ExpandStringSlicePtr(input []interface{}) *[]string {
	result := ExpandStringSlice(input)
	return &result
}

func FlattenStringSlice(input []string) []interface{} {
	result := make([]interface{}, 0)
	for _, item := range input {
		result = append(result, item)
	}
	return result
}

func FlattenStringSlicePtr(input *[]string) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		result = FlattenStringSlice(*input)
	}
	return result
}
