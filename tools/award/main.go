// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"os"
	"strconv"
)

// 恒星年出块数
// 6分钟 87661
// 8分钟 65746
const YEARBLOCKS = 65746

var yNext int

// AwardList 计算减量发行总量。
// @base: 初始币量/块
// @rate: 递减率
// @ny: 步进年数（固定币量持续）
// @stop: 块币停止值
// @return: 总量
func AwardList(base, rate, ny, stop int) int {
	sum := 0
	y := yNext
	// 到每块3币时停止。
	for base >= stop {
		ysum := base * YEARBLOCKS * ny
		sum += ysum
		y += ny
		fmt.Printf("%d\t%-8d\t[%-7d]\t%d\n", y, sum, ysum, base)

		base = base * rate / 100
	}
	return sum
}

// Award3y 计算初期n年的总量。
func Award3y(n int) int {
	sum := 0
	ysum := 0
	base := 0

	for y:=1; y<=n; y++ {
		base = y * 10
		ysum = YEARBLOCKS * base
		sum += ysum
		fmt.Printf("%d\t%-8d\t[%-7d]\t%d\n", y, sum, ysum, base)
	}
	yNext = n

	return sum
}

func main() {
	// 从命令行参数获取基础币量和年利率。
	// 例：go run main.go 40 80 2 3? // 末尾的 ? 表示可选
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <base> <rate>")
		return
	}

	base, err1 := strconv.Atoi(os.Args[1])
	rate, err2 := strconv.Atoi(os.Args[2])
	ny, err3 := strconv.Atoi(os.Args[3])
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("Invalid input. Please provide two integers.")
		return
	}
	// 默认3币/块后终止
	stop := 3
	if len(os.Args) == 5 {
		stop, err1 = strconv.Atoi(os.Args[4])
		if err1 != nil {
			fmt.Println("Invalid input. Please provide two integers.")
			return
		}
	}

	fmt.Println("年次\t累计\t\t（次计）\t币量/块")
	fmt.Println("-------------------------------------------------------")

	test3 := Award3y(3)
	fmt.Println("-------------------------------------------------------")

	total := AwardList(base, rate, ny, stop)

	fmt.Println("-------------------------------------------------------")
	fmt.Println("初期三年：", test3)
	fmt.Println("减量发行：", total)
	fmt.Println("发行总计：", total+test3)
}
