package utils

func AddToSlice(slice *[]string, value string) {
	*slice = append(*slice, value)
}

func RemoveFromSliceByValue(slice *[]string, value string) {
	index := -1
	for i, v := range *slice {
		if v == value {
			index = i
			break
		}
	}

	if index != -1 {
		*slice = append((*slice)[:index], (*slice)[index+1:]...)
	}
}

func RemoveFromSlice(slice []map[string]string, s int) []map[string]string {
	return append(slice[:s], slice[s+1:]...)
}
