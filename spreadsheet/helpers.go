package spreadsheet

func indexOf(in []string, target string) int {
	for index, header := range in {
		if header == target {
			return index
		}
	}
	return -1
}
