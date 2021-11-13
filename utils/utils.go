package utils

type IndexPred func(i int) bool

// 扫描slice是否包含符合条件的项
func FindIf(lenOfSlice int, pred IndexPred) int {
	for i := 0; i < lenOfSlice; i++ {
		if pred(i) {
			return i
		}
	}
	return -1
}
