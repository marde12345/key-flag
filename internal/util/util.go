package util

func ConvertArrayInterfaceToArrayString(input []interface{}) []string {
	result := make([]string, 0)
	for _, row := range input {
		result = append(result, row.(string))
	}

	return result
}
