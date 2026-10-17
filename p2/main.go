package main

import "strconv"
import "fmt"

// digits 返回房间号 x 的十进制数字。
func digits(x int) []byte {
	return []byte(strconv.Itoa(x))
}

// reverse 逆转数字序列,忽略前导零。
func reverse(d []byte) int {
	// 逆转后末尾的零对应原数的前导零,直接忽略。
	i := 0
	for i < len(d) && d[i] == '0' {
		i++
	}
	if i == len(d) {
		return 0
	}
	d = d[i:]
	for l, r := 0, len(d)-1; l < r; l, r = l+1, r-1 {
		d[l], d[r] = d[r], d[l]
	}
	y, err := strconv.Atoi(string(d))
	if err != nil {
		return -1
	}
	return y
}

// state 表示一个 (房间, 模式) 状态。
type state struct {
	room, mode int
}

func main() {
	var M, N int
	fmt.Scan(&M, &N)

	// 抵达房间 N 后即可结束任务,无需切换到特定模式。
	visited := make([][3]bool, M+1)
	visited[1][0] = true
	frontier := []state{{1, 0}}
	for steps := 0; ; steps++ {
		var next []state
		for _, x := range frontier {
			if x.room == N {
				fmt.Println(steps)
				return
			}
			switch x.mode {
			case 0:
				// 模式一:走廊, x -> x+1。
				if x.room+1 <= M && !visited[x.room+1][0] {
					visited[x.room+1][0] = true
					next = append(next, state{x.room + 1, 0})
				}
			case 1:
				// 模式二:跃迁, x -> 2x。
				if 2*x.room <= M && !visited[2*x.room][1] {
					visited[2*x.room][1] = true
					next = append(next, state{2 * x.room, 1})
				}
				// 模式二:跃迁, x -> x/2,仅当 x 为偶数时可整除。
				if x.room%2 == 0 && x.room/2 >= 1 && !visited[x.room/2][1] {
					visited[x.room/2][1] = true
					next = append(next, state{x.room / 2, 1})
				}
			case 2:
				// 模式三:逆转(忽略前导零),产生的房间号必须在 [1, M] 内。
				y := reverse(digits(x.room))
				if y >= 1 && y <= M && !visited[y][2] {
					visited[y][2] = true
					next = append(next, state{y, 2})
				}
			}
			// 模式切换:切换到另外两种模式之一,房间不变,计 1 步。
			for _, m := range []int{0, 1, 2} {
				if m != x.mode && !visited[x.room][m] {
					visited[x.room][m] = true
					next = append(next, state{x.room, m})
				}
			}
		}
		frontier = next
	}
}
