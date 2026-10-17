package main

import "fmt"

// gcd 返回 a 与 b 的最大公约数。
func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	var p, q int64
	fmt.Scan(&p, &q)

	// 约分到最简,数值相等的房间中编号最小的即最简分数首次出现的位置。
	g := gcd(p, q)
	p /= g
	q /= g

	// 所在斜线编号: 斜线 k 包含全部满足 p+q-1=k 的分数,共 k 个。
	k := p + q - 1
	// 斜线 k 之前的房间总数。
	idx := k * (k - 1) / 2
	// 偶数斜线 p 递增,奇数斜线 p 递减。
	if k%2 == 0 {
		idx += p - 1
	} else {
		idx += q - 1
	}

	// 每排两间房,经过的排数即房间编号整除 2。
	fmt.Println(idx / 2)
}
