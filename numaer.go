package numaer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

// Node ： NUMA Node 节点信息
type Node struct {
	Name string
}

// Zone NUMA  zone 区域信息
type Zone struct {
	Type string
	FreePage int
	Node *Node
}

// CPU 相关信息
type CPU struct {
	//core id	
	ID int
	Node *Node
}

// IsNUMA  判断是否为 NUMA 架构
func IsNUMA() bool {
	if _, err := os.Stat("/proc/zoneinfo"); !os.IsNotExist(err) {
		return true
	} 
	return false
}

// Nodes 获取当前系统 内存节点 node 信息
func Nodes() ([]*Node, error) {
	if !IsNUMA() {
		return nil, fmt.Errorf("OS is not NUMA")
	}

	f, err := os.Open("/proc/buddyinfo")
	// buddy 文件包含了node 相关信息
	// 如：
	/*
	Node 0, zone      DMA     29      9      9      6      2      2      1      1      2      2      0 
	Node 0, zone    DMA32    941   2242   1419    398    142     60     16      1      1      0      0 
	*/

	if err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}
	defer f.Close()

	var NUMANodeSlice []string

	// NewScanner创建并返回一个从f读取数据的Scanner，默认的分割函数是ScanLines
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, ",") // 以,切片
		NUMANodeSlice = append(NUMANodeSlice, fields[0])
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}

	// 去重
	NUMANodeSlice = removeReplicaSliceString(NUMANodeSlice)

	var Nodes []*Node
	
	for _, NodeName := range NUMANodeSlice {
		Nodes = append(Nodes, &Node{
			Name: NodeName,
		})
	}

	return Nodes, nil
} 

// NumNode ：获取当前系统 NUMA 数量
func NumNode() (int, error) {

	// 获取当前系统参数
	if NodeSlice, err := Nodes(); err == nil {
		return len(NodeSlice), nil  //返回 []Node 数量
	} else {
		return 0, err
	}
}


// ZoneInfo 获取内存节点 node 的区域信息
func ZoneInfo(n *Node) ([]*Zone, error) {

	f, err := os.Open("/proc/zoneinfo")
	// zoneinfo 文件包含了zone 相关信息
	// 如：
	/*
	Node 0, zone      DMA
	 pages free     3969 
	...
	  pagesets
        cpu: 0
              count: 0
              high:  0
              batch: 1
  		vm stats threshold: 8
    	cpu: 1
              count: 0
              high:  0
              batch: 1
  		vm stats threshold: 8
    	cpu: 2
              count: 0
              high:  0
              batch: 1
  		vm stats threshold: 8
    	cpu: 3
              count: 0
              high:  0
              batch: 1
	*/

	if err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}
	defer f.Close()

	var ZoneSlice []*Zone

	pageTag := false
	var tmpZoneType string

	// NewScanner创建并返回一个从f读取数据的Scanner，默认的分割函数是ScanLines
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		txt = strings.TrimSpace(txt) //去除首尾空格
		fields := strings.Split(txt, " ") // 以  空格 切片 Node 0, zone DMA
		fields = removeNullSliceString(fields)

		// 关于 page free 的条目一般都在 Zone 信息后一排
		// 在上一行信息中获取到了 Zone 信息，保存 temZoneType 中，并设置 pageTag = true
		if pageTag && len(fields) >= 3 && fields[0] == "pages" && fields[1] == "free" {   //pages free 969

			pagefree, err := strconv.Atoi(fields[2])

			// 异常则设置为0
			if err != nil {
				pagefree = 0
			}

			ZoneSlice = append(ZoneSlice, &Zone{
				Type: tmpZoneType,
				Node: n,
				FreePage: pagefree, 
			})
			// 设置  pageTag = false
			pageTag = false
			// 跳过本次循环
			continue
		}

		// Node 0, zone DMA 
		if len(fields) >= 4 && fields[0] + fields[1] == n.Name + "," && fields[2] == "zone" {
			tmpZoneType = fields[3]
			pageTag = true
		}
		
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}
	return ZoneSlice, nil
}

// CPUInfo 获取内存节点 Node 绑定的CPU信息
func CPUInfo(n *Node) ([]*CPU, error) {
	return nil,nil
}

// BuddyInfo ：伙伴系统当前状态
func BuddyInfo(z *Zone) (map[int]int64, error) {// [11中内存碎片大小]剩余碎片数

	NodeName := z.Node.Name
	ZoneType := z.Type

	f, err := os.Open("/proc/buddyinfo")
	// buddy 文件包含了node 相关信息
	// 如：
	/*
	Node 0, zone      DMA     29      9      9      6      2      2      1      1      2      2      0 
	Node 0, zone    DMA32    941   2242   1419    398    142     60     16      1      1      0      0 
	*/

	if err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}

	defer f.Close()

	buddyMap := make(map[int]int64)

	// NewScanner创建并返回一个从f读取数据的Scanner，默认的分割函数是ScanLines
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		buddySlice := removeNullSliceString(strings.Split(txt, " "))
		// 判断相关信息
		if (buddySlice[0] + buddySlice[1]) == (NodeName + ",") && buddySlice[2] == "zone" && buddySlice[3] == ZoneType {
			for index, v := range buddySlice[4:] {
				vv, _ := strconv.ParseInt(v, 10, 64)
				k := 0
				if index == 0 {
					buddyMap[k] = vv
				} else {
					// 开方  math库的要f64 太难了转化了
					for i := 0; i < index; i++ {
						k = k * 2
					}
					buddyMap[k] = vv
				}
			}
		}  
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}	
	return buddyMap, nil
}