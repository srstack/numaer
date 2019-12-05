package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

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

// Node ： NUMA Node 节点信息
type Node struct {
	NodeID int		`json:"node ID"`
	Name string   	`json:"name"`
}

// Zone ： NUMA  zone 区域信息
type Zone struct {
	Type string		`json:"type"`
	FreePage int	`json:"free page"`
	Node *Node   	`json:"node"`
}

// CPU ：相关信息
type CPU struct {
	//core id	
	CoreID int 					`json:"core id"`
	Node *Node 					`json:"node"`
	CPUInfo map[string]string	`json:"cpu info"`
}

// IsNUMA ：判断是否为 NUMA 架构
func IsNUMA() bool {
	if _, err := os.Stat("/proc/zoneinfo"); !os.IsNotExist(err) {
		return true
	} 
	return false
}

// GetCPUInfo : 通过id获取CPU信息
func GetCPUInfo(coreID int) (*CPU, error) {

	if !IsNUMA() {
		return nil, fmt.Errorf("OS is not NUMA")
	}

	f, err := os.Open("/proc/cpuinfo")

	/*
	processor       : 0
	vendor_id       : GenuineIntel
	cpu family      : 6
	model           : 79
	model name      : Intel(R) Xeon(R) CPU E5-26xx v4
	stepping        : 1
	microcode       : 0x1
	cpu MHz         : 2394.446
	cache size      : 4096 KB
	physical id     : 0
	siblings        : 1
	core id         : 0
	cpu cores       : 1
	apicid          : 0
	initial apicid  : 0
	fpu             : yes
	fpu_exception   : yes
	cpuid level     : 13
	wp              : yes
	flags           : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch bmi1 avx2 bmi2 rdseed adx xsaveopt
	bogomips        : 4788.89
	clflush size    : 64
	cache_alignment : 64
	address sizes   : 40 bits physical, 48 bits virtual
	power management:
	*/

	if err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}
	defer f.Close()

	cpuinfo := &CPU{
		CoreID: coreID,	
		CPUInfo: make(map[string]string),
	}

	findIDTag := false

	// NewScanner创建并返回一个从f读取数据的Scanner，默认的分割函数是ScanLines
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, ":") // 以:切片
		fields = removeNullSliceString(fields)

		// 目标 cpu 信息获取完毕
		if findIDTag && len(fields) > 0 && strings.TrimSpace(fields[0]) == "processor" {
			break
		}

		// 找到目标 cpu
		if len(fields) == 2 && strings.TrimSpace(fields[0]) == "processor" && strings.TrimSpace(fields[1]) == strconv.Itoa(cpuinfo.CoreID) {
			findIDTag = true
		}

		// 存入信息
		if findIDTag && len(fields) > 1 {
			cpuinfo.CPUInfo[strings.TrimSpace(fields[0])] = strings.TrimSpace(fields[1])
		}
		if len(fields) == 1 && findIDTag {
			cpuinfo.CPUInfo[strings.TrimSpace(fields[0])] = ""
		}
	}

	if !findIDTag {
		return nil, fmt.Errorf("err: can't find cpu %v", coreID)
	}

	// 获取该 cpu 的 node 信息
	cpuinfo, err = cpuinfo.CPUGetNodeInfo()

	if err != nil {
		return cpuinfo, fmt.Errorf("err: can't get node info %v", err)
	}

	return cpuinfo, nil
}

// CPUGetNodeInfo : 通过cpu core id 找到对应的内存 node 节点信息
func (c *CPU) CPUGetNodeInfo() (*CPU, error) {
	f, err := os.Open("/proc/zoneinfo")
	/*
	Node 0, zone      DMA
	pages free     3969
			min      3
			low      3
			high     4
			scanned  0
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
	...
	*/

	if err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}
	defer f.Close()

	var nodeName string

	// NewScanner创建并返回一个从f读取数据的Scanner，默认的分割函数是ScanLines
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ") // 以,切片
		fields = removeNullSliceString(fields)

		// Node 0, zone DMA 
		if len(fields) == 4 && fields[0] == "Node" && fields[2] == "zone" {
			// join 拼接
			nodeName = strings.Join([]string{fields[0], fields[1]}, "") // node0,
		}

		// 匹配对应的 core id ， 则上一次匹配的 node name 为所需 node name
		if nodeName != "" && fields[0] == "cpu:" && fields[1] == strconv.Itoa(c.CoreID) {
			break
		}	
	}

	// Node1,  => 1
	nodeID, err := strconv.Atoi(substr(nodeName,len(nodeName)-2, 1))

	if err != nil {
		return nil, fmt.Errorf("err: str >> int err %v", err)
	} 

	// 根据 获取 node info
	nodeinfo, err := GetNodeInfo(nodeID)

	if err != nil {
		return nil, fmt.Errorf("err: get node info err %v", err)
	}

	c.Node = nodeinfo

	return c, nil
}

// Nodes ：获取当前系统 内存节点 node 信息
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
		NUMANodeSlice = append(NUMANodeSlice, strings.Replace(fields[0], " ", "", -1)) // Node 0 => Node0
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}

	// 去重
	NUMANodeSlice = removeReplicaSliceString(NUMANodeSlice)

	var Nodes []*Node
	
	for _, NodeName := range NUMANodeSlice {
		id, err := strconv.Atoi(substr(NodeName, len(NodeName)-1, 1))
		if err != nil {
			return nil, fmt.Errorf("err: str >> int err %v", err)
		} 
		Nodes = append(Nodes, &Node{
			NodeID: id,
			Name: NodeName,
		})
	}

	return Nodes, nil
} 

// GetNodeInfo ：通过 node 节点 id 获取 node 信息
func GetNodeInfo(nodeID int) (*Node, error) {

	if !IsNUMA() {
		return nil, fmt.Errorf("OS is not NUMA")
	}

	// 获取全部node 信息
	nodes, err := Nodes()
	
	if err != nil {
		return nil, fmt.Errorf("error : Get nodes fail %v", err)
	}

	for _, node := range nodes {
		if node.NodeID == nodeID {
			return node, nil
		}
	}

	return nil, fmt.Errorf("error: can't find node %v", nodeID)
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


// ZoneInfo : 获取内存节点 node 的 zone 区域信息
func (n *Node) ZoneInfo() ([]*Zone, error) {

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
		if pageTag && fields[0] == "pages" && len(fields) >= 3  && fields[1] == "free" {   //pages free 969

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
		if len(fields) == 4 && strings.Join([]string{fields[0], fields[1]}, "") == strings.Join([]string{n.Name, ","}, "") && fields[2] == "zone" {
			tmpZoneType = fields[3]
			pageTag = true
		}
		
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}
	return ZoneSlice, nil
}

// CPUInfo ：获取内存节点 Node 绑定的CPU信息
func (n *Node) CPUInfo() ([]*CPU, error) {
	f, err := os.Open("/proc/zoneinfo")
	/*
	Node 0, zone      DMA
	pages free     3969
			min      3
			low      3
			high     4
			scanned  0
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
	...
	*/

	if err != nil {
		return nil, fmt.Errorf("err : %v", err)
	}
	defer f.Close()

	matchNode := false
	var cupIDStringSlice []string

	// NewScanner创建并返回一个从f读取数据的Scanner，默认的分割函数是ScanLines
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ") // 以,切片
		fields = removeNullSliceString(fields)

		// Node 0, zone DMA 
		if len(fields) == 4 && fields[0] == "Node" && strings.Join([]string{fields[0], fields[1]}, "") == strings.Join([]string{n.Name, ","}, "") {
			// 第一次匹配该节点后 设置 matchNode 为 ture
			// 若 matchNode 已经为 ture 则说明匹配完成，直接退出
			if matchNode {
				break
			} else {
				matchNode = true
			}
		}

		if matchNode && fields[0] == "cpu:" {
			cupIDStringSlice = append(cupIDStringSlice, fields[1])
		}
	}
	
	var CPUSlice []*CPU

	// 去重
	cupIDStringSlice = removeReplicaSliceString(cupIDStringSlice)

	// 完成 CUP slice
	for _, cpuID := range cupIDStringSlice {
		ID, err := strconv.Atoi(cpuID)
		if err != nil {
			return nil, fmt.Errorf("err: str >> int err %v", err)
		} 

		// 通过 id 获取 cpu info
		cpuinfo, err := GetCPUInfo(ID)
		if err != nil {
			return nil, fmt.Errorf("err: get cpu info err %v", err)
		}

		CPUSlice = append(CPUSlice, cpuinfo)
	}

	return CPUSlice, nil
}

// BuddyInfo ：伙伴系统当前状态
func (z *Zone) BuddyInfo() (map[int]int64, error) {// [11中内存碎片大小]剩余碎片数

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
		if strings.Join([]string{buddySlice[0], buddySlice[1]}, "") == strings.Join([]string{NodeName, ","}, "") && buddySlice[2] == "zone" && buddySlice[3] == ZoneType {
			for index, v := range buddySlice[4:] {
				vv, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					return nil, fmt.Errorf("err: str >> int err %v", err)
				} 
				k := 1
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

func main(){
	if IsNUMA() {
		fmt.Println("os is NUMA")
	}

	Nodes, err := Nodes()

	if err != nil {
		fmt.Errorf("ERR: %v" , err)
	}

	fmt.Print("--------------------\n")
	for _, v := range Nodes {
		fmt.Printf("Node ID: %v,\nNode Name: %v \n",v.NodeID , v.Name)
	}
	fmt.Print("--------------------\n")


	if numNode, err := NumNode(); err == nil {
		fmt.Printf("NUMA Node Num: %v \n",numNode)
	}
	fmt.Print("--------------------\n")

	Zones, err := Nodes[0].ZoneInfo()

	if err != nil {
		fmt.Errorf("ERR: %v" , err)
	}


	for i, v := range Zones {
		fmt.Printf("Zone: %v \n",i)
		fmt.Printf("Zone Type: %v \n",v.Type)
		fmt.Printf("Zone PageFree: %v \n",v.FreePage)
	}
	fmt.Print("--------------------\n")

	buddy,err := Zones[0].BuddyInfo()

	fmt.Printf("Node: %v ,\nZone: %v,\n" , Nodes[0].Name, Zones[0].Type)

	fmt.Print("--------------------\n")

	fmt.Print("BuddyInfo: \n")
	for k, v := range buddy { 
		fmt.Printf("page block size: %v, free num: %v \n",k,v)
	}
	fmt.Print("--------------------\n")

	CPUs ,err := Nodes[0].CPUInfo()

	if err != nil {
		fmt.Errorf("ERR: %v" , err)
	}

	for _, v := range CPUs {
		fmt.Printf("CPU ID : %v,\nCPU Node ID:%v\n", v.CoreID, v.Node.NodeID)
	}

	fmt.Print("--------------------\n")

	fmt.Printf("CPU %v Info: \n", CPUs[0].CoreID)

	for k, v := range CPUs[0].CPUInfo { 
		fmt.Printf("%v:%v\n",k , v)
	}
	fmt.Print("--------------------\n")

	cpu1, err := GetCPUInfo(0)

	if err != nil {
		fmt.Errorf("ERR: %v" , err)
	}

	fmt.Printf("Get Cpu 1: %v\n", cpu1) 

	fmt.Print("--------------------\n")

	node1, err := GetNodeInfo(0)

	if err != nil {
		fmt.Errorf("ERR: %v" , err)
	}

	fmt.Printf("Get node 1: %v\n", node1) 

}

