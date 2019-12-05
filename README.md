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

### test output

    os is NUMA
    --------------------
    Node ID: 0,
    Node Name: Node0 
    --------------------
    NUMA Node Num: 1 
    --------------------
    Zone: 0 
    Zone Type: DMA 
    Zone PageFree: 1936 
    Zone: 1 
    Zone Type: DMA32 
    Zone PageFree: 23421 
    --------------------
    Node: Node0 ,
    Zone: DMA,
    --------------------
    BuddyInfo: 
    page block size: 32, free num: 8 
    page block size: 64, free num: 4 
    page block size: 256, free num: 2 
    page block size: 4, free num: 12 
    page block size: 2, free num: 4 
    page block size: 8, free num: 11 
    page block size: 16, free num: 7 
    page block size: 128, free num: 1 
    page block size: 512, free num: 1 
    page block size: 1024, free num: 0 
    page block size: 1, free num: 16 
    --------------------
    CPU ID : 0,
    CPU Node ID:0
    --------------------
    CPU 0 Info: 
    cpuid level:13
    processor:0
    vendor_id:GenuineIntel
    cpu family:6
    stepping:1
    microcode:0x1
    siblings:1
    core id:0
    cache_alignment:64
    physical id:0
    apicid:0
    wp:yes
    power management:
    cache size:4096 KB
    cpu cores:1
    initial apicid:0
    fpu:yes
    flags:fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch bmi1 avx2 bmi2 rdseed adx xsaveopt
    bogomips:4788.89
    model:79
    model name:Intel(R) Xeon(R) CPU E5-26xx v4
    cpu MHz:2394.446
    fpu_exception:yes
    clflush size:64
    address sizes:40 bits physical, 48 bits virtual
    --------------------
    Get Cpu 1: &{0 0xc000099820 map[address sizes:40 bits physical, 48 bits virtual apicid:0 bogomips:4788.89 cache size:4096 KB cache_alignment:64 clflush size:64 core id:0 cpu MHz:2394.446 cpu cores:1 cpu family:6 cpuid level:13 flags:fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ss ht syscall nx lm constant_tsc rep_good nopl eagerfpu pni pclmulqdq ssse3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch bmi1 avx2 bmi2 rdseed adx xsaveopt fpu:yes fpu_exception:yes initial apicid:0 microcode:0x1 model:79 model name:Intel(R) Xeon(R) CPU E5-26xx v4 physical id:0 power management: processor:0 siblings:1 stepping:1 vendor_id:GenuineIntel wp:yes]}
    --------------------
    Get node 1: &{0 Node0}