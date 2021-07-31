package zcomplier

var strTable []string

// 返回新增str在symbolList中的索引
func addstr(str string) int {
	for i := 0; i < len(strTable); i++ {
		s := strTable[i]
		if s == str {
			return i
		}
	}
	strTable = append(strTable, str)
	return len(strTable) - 1
}

func getStrByIndex(index int) string {
	if index < len(strTable) {
		return strTable[index]
	}
	return ""
}
