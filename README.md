# numaer
 Go操作NUMA相关基础库

### test code

    func main(){
        if IsNUMA() {
            fmt.Println("os is NUMA")
        }

        Nodes, err := Nodes()

        if err != nil {
            fmt.Errorf("ERR: %v" , err)
        }

        for _, v := range Nodes {
            fmt.Printf("node : %v \n",v.Name )
        }

        if numNode, err := NumNode(); err == nil {
            fmt.Printf("NUMA Node Num: %v \n",numNode)
        }

        Zones, err := Nodes[0].ZoneInfo()

        if err != nil {
            fmt.Errorf("ERR: %v" , err)
        }

        for i, v := range Zones {
            fmt.Printf("Zone: %v \n",i)
            fmt.Printf("Zone Type: %v \n",v.Type)
            fmt.Printf("Zone PageFree: %v \n",v.FreePage)
        }

        buddy,err := Zones[0].BuddyInfo()

        fmt.Printf("Node: %v ,Zone: %v , buddyInfo:\n" , Nodes[0].Name, Zones[0].Type)

        for k, v := range buddy { 
            fmt.Printf("page block size: %v, free num: %v \n",k,v)
        }

    }

### test output

    os is NUMA
    node : Node0 
    NUMA Node Num: 1 
    Zone: 0 
    Zone Type: DMA 
    Zone PageFree: 1938 
    Zone: 1 
    Zone Type: DMA32 
    Zone PageFree: 26748 
    Node: Node0 ,Zone: DMA , buddyInfo:
    page block size: 2, free num: 16 
    page block size: 4, free num: 9 
    page block size: 128, free num: 1 
    page block size: 512, free num: 2 
    page block size: 1, free num: 14 
    page block size: 8, free num: 6 
    page block size: 16, free num: 3 
    page block size: 32, free num: 1 
    page block size: 64, free num: 1
    page block size: 2, free num: 16 
    page block size: 16, free num: 3 
    page block size: 128, free num: 1 
    page block size: 1024, free num: 0 