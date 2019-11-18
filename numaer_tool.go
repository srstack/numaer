package numaer


// RemoveReplicaSliceString 切片去重
func RemoveReplicaSliceString(srcSlice []string) []string {
 
	resultSlice := make([]string, 0)
	// 利用map key 值唯一去重
    tempMap := make(map[string]bool, len(srcSlice))
    for _, v := range srcSlice{
        if tempMap[v] == false{
            tempMap[v] = true
            resultSlice = append(resultSlice, v)
        }
    }
    return resultSlice
}


// RemoveNullSliceString : 删除空白字符的元素
func RemoveNullSliceString(srcSlice []string) []string {
 
	resultSlice := make([]string, 0)

	// 循环判断
    for _, v := range srcSlice{
        if v != "" && v != " " {
            resultSlice = append(resultSlice, v)
        }
    }
    return resultSlice
}