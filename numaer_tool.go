package numaer


// RemoveReplicaSliceString : 切片去重
func removeReplicaSliceString(srcSlice []string) []string {
 
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
func removeNullSliceString(srcSlice []string) []string {
 
	resultSlice := make([]string, 0)

	// 循环判断
    for _, v := range srcSlice{
        if v != "" && v != " " {
            resultSlice = append(resultSlice, v)
        }
    }
    return resultSlice
}

// susbstr ：字符串截取
func substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0

    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length

    if start > end {
        start, end = end, start
    }

    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }

    return string(rs[start:end])
}