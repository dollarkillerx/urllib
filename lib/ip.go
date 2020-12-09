package lib

import (
	"fmt"
	"math/rand"
	"time"
)

// 随机生成合法 IP，如： 222.16.123.95
func RandomIp() string {
	// IP 范围二维数组
	ranges := ipRange()
	idx := newRand().Intn(10)
	return numToIp(ranges[idx][0] + newRand().Intn(ranges[idx][1]-ranges[idx][0]))
}

// 随机生成（隐蔽后两位的）合法 IP，如： 222.16.*.*
func RandomOmicIp() string {
	// IP 范围二维数组
	ranges := ipRange()
	idx := newRand().Intn(10)
	return numToOmicIp(ranges[idx][0] + newRand().Intn(ranges[idx][1]-ranges[idx][0]))
}

func numToIp(num int) string {
	var arr []int = make([]int, 4)
	arr[0] = (num >> 24) & 0xff
	arr[1] = (num >> 16) & 0xff
	arr[2] = (num >> 8) & 0xff
	arr[3] = num & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

func numToOmicIp(num int) string {
	var arr []int = make([]int, 2)
	arr[0] = (num >> 24) & 0xff
	arr[1] = (num >> 16) & 0xff
	return fmt.Sprintf("%d.%d.*.*", arr[0], arr[1])
}

// IP 范围二维数组
func ipRange() [][]int {
	return [][]int{{607649792, 608174079}, //36.56.0.0-36.63.255.255
		{1038614528, 1039007743},   //61.232.0.0-61.237.255.255
		{1783627776, 1784676351},   //106.80.0.0-106.95.255.255
		{2035023872, 2035154943},   //121.76.0.0-121.77.255.255
		{2078801920, 2079064063},   //123.232.0.0-123.235.255.255
		{-1950089216, -1948778497}, //139.196.0.0-139.215.255.255
		{-1425539072, -1425014785}, //171.8.0.0-171.15.255.255
		{-1236271104, -1235419137}, //182.80.0.0-182.92.255.255
		{-770113536, -768606209},   //210.25.0.0-210.47.255.255
		{-569376768, -564133889},   //222.16.0.0-222.95.255.255
	}
}

// 实例化随机数结构体，源为时间微秒
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
