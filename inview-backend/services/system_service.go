package services

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"insense-local/config"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type SystemServiceInterface interface {
	SystemInfo() map[string]interface{}
	SystemLog(env *config.Env) (string, error)
	NetworkInfo() (map[string]string, error)
}

func SystemService(timeout time.Duration) SystemServiceInterface {
	return &systemService{
		contextTimeout: timeout,
	}
}

type systemService struct {
	contextTimeout time.Duration
}

func formatBytesToGB(bytes uint64) string {
	return fmt.Sprintf("%.2f", float64(bytes)/(1024*1024*1024))
}

func formatTime(t uint64) string {
	return time.Unix(int64(t), 0).Format("2006-01-02 15:04:05")
}
func (ss *systemService) SystemInfo() map[string]interface{} {
	statMap := make(map[string]interface{})

	// 获取主机相关信息
	hostInfo, _ := host.Info()
	hostMap := make(map[string]interface{})
	hostMap["uptime"] = hostInfo.Uptime                   // 运行时间
	hostMap["bootTime"] = formatTime(hostInfo.BootTime)   // 启动时间
	hostMap["procs"] = hostInfo.Procs                     // 进程数
	hostMap["os"] = hostInfo.OS                           // 操作系统
	hostMap["platform"] = hostInfo.Platform               // 平台
	hostMap["platformVersion"] = hostInfo.PlatformVersion // 平台版本
	hostMap["kernelArch"] = hostInfo.KernelArch           // 内核
	hostMap["kernelVersion"] = hostInfo.KernelVersion     // 内核版本
	statMap["hosts"] = hostMap

	// 获取内存信息
	memInfo, _ := mem.VirtualMemory()
	memMap := make(map[string]interface{})
	memMap["total"] = formatBytesToGB(memInfo.Total)         // 总内存
	memMap["available"] = formatBytesToGB(memInfo.Available) // 可用内存
	memMap["used"] = formatBytesToGB(memInfo.Used)           // 已使用内存
	memMap["free"] = formatBytesToGB(memInfo.Free)           // 剩余内存
	memMap["usedPercent"] = memInfo.UsedPercent              // 百分比
	memMap["buffers"] = formatBytesToGB(memInfo.Buffers)     // 缓存
	memMap["shared"] = formatBytesToGB(memInfo.Shared)       // 共享内存
	memMap["cached"] = formatBytesToGB(memInfo.Cached)       // 缓冲区
	statMap["mems"] = memMap

	// 获取CPU信息
	//cpuInfo, _ := cpu.Info()
	//var cpuMapArr []map[string]interface{}
	//for _, c := range cpuInfo {
	//	cpuMap := make(map[string]interface{})
	//	cpuMap["cpu"] = c.CPU + 1         // 第几个CPU 从0开始的
	//	cpuMap["cores"] = c.Cores         // CPU的核数
	//	cpuMap["modelName"] = c.ModelName // CPU类型
	//	cpuMapArr = append(cpuMapArr, cpuMap)
	//}
	//statMap["cpus"] = cpuMapArr

	// 获取IO信息
	//ioInfo, _ := net.IOCounters(false)
	//var ioMapArr []map[string]interface{}
	//for _, i := range ioInfo {
	//	ioMap := make(map[string]interface{})
	//	ioMap["ioName"] = i.Name             // 网口名
	//	ioMap["bytesSent"] = i.BytesSent     // 发送字节数
	//	ioMap["bytesRecv"] = i.BytesRecv     // 接收字节数
	//	ioMap["packetsSent"] = i.PacketsSent // 发送的数据包数
	//	ioMap["packetsRecv"] = i.PacketsRecv // 接收的数据包数
	//	ioMapArr = append(ioMapArr, ioMap)
	//}
	//statMap["ios"] = ioMapArr

	// 获取磁盘信息
	partitions, _ := disk.Partitions(false)
	var diskMapArr []map[string]interface{}
	for _, partition := range partitions {
		diskMap := make(map[string]interface{})
		usage, _ := disk.Usage(partition.Mountpoint)
		diskMap["disk"] = partition.Mountpoint          // 第几块磁盘
		diskMap["total"] = formatBytesToGB(usage.Total) // 总大小
		diskMap["free"] = formatBytesToGB(usage.Free)   // 剩余空间
		diskMap["used"] = formatBytesToGB(usage.Used)   // 已使用空间
		diskMap["usedPercent"] = usage.UsedPercent      // 百分比
		diskMapArr = append(diskMapArr, diskMap)
	}
	statMap["disks"] = diskMapArr

	return statMap
}

func (ss *systemService) SystemLog(env *config.Env) (string, error) {
	// 读取日志文件内容
	file, err := os.Open(env.LogPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (ss *systemService) NetworkInfo() (map[string]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	fmt.Println(interfaces)

	for _, iface := range interfaces {
		// 过滤掉虚拟网络接口和回环接口
		if strings.HasPrefix(iface.Name, "vEthernet") ||
			strings.HasPrefix(iface.Name, "virbr") ||
			strings.HasPrefix(iface.Name, "docker") ||
			strings.HasPrefix(iface.Name, "lo") {
			continue
		}

		// 获取接口的所有 IP 地址
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {

			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ip = ip.To4()
			// 跳过无效IP
			if ip == nil || ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() {
				continue
			}

			// 跳过 APIPA 地址
			if ip.IsLinkLocalUnicast() && strings.HasPrefix(ip.String(), "169.254.") {
				continue
			}

			// 构建结果并返回
			info := make(map[string]string)
			info["interfaceName"] = iface.Name
			info["ip"] = ip.String()

			// 获取对应的 MAC 地址，并格式化为 machineId
			mac := iface.HardwareAddr.String()
			if mac != "" {
				machineId := strings.ReplaceAll(mac, ":", "")
				info["machineID"] = machineId
				info["mac"] = mac
			} else {
				info["mac"] = "N/A"
				info["machineID"] = "N/A"
			}

			// 找到一个能上网的接口，直接返回
			return info, nil
		}
	}

	return nil, fmt.Errorf("未找到可以上网的网络接口")
}
