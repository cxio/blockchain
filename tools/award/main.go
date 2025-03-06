// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"os"
	"strconv"
)

// 恒星年出块数（6分钟）。
const SY6BLOCKS = 87661

// AwardList 计算减量发行总量。
// @base: 初始币量/块
// @rate: 递减率
// @ny: 步进年数（固定币量持续）
// @stop: 块币停止值
// @return: 总量
func AwardList(base, rate, ny, stop int) int {
	base *= 1

	fmt.Println("年次\t累计\t\t（次计）\t币量/块")
	fmt.Println("--------------------------------------------------------")

	sum := 0
	y := 0
	// 到每块3币时停止。
	for base >= stop {
		ysum := base * SY6BLOCKS * ny
		sum += ysum
		y += ny
		fmt.Printf("%d\t%d \t(%d)\t%d\n", y, sum, ysum, base)

		base = base * rate / 100
	}
	return sum
}

// Award3y 计算初期三年的总量。
func Award3y() int {
	return SY6BLOCKS*10 + SY6BLOCKS*20 + SY6BLOCKS*30
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
	test3 := Award3y()
	total := AwardList(base, rate, ny, stop)

	fmt.Println("--------------------------------------------------------")
	fmt.Println("初期三年：", test3)
	fmt.Println("减量发行：", total)
	fmt.Println("发行总计：", total+test3)
}
