package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// env 读取环境变量,未设置或为空时返回默认值。
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// jwtRe 匹配 JWT 形状的通行证。
var jwtRe = regexp.MustCompile(`eyJ[A-Za-z0-9_-]*\.eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]+`)

// firstJWT 返回文本中第一个 JWT。
func firstJWT(s string) string {
	return jwtRe.FindString(s)
}

// postSeq 向 /api/challenge 提交一个编号。
func postSeq(httpc *http.Client, base, jwt string, seq int) string {
	req, err := http.NewRequest(http.MethodPost, base+"/api/challenge", strings.NewReader(fmt.Sprintf(`{"seq": %d}`, seq)))
	if err != nil {
		return err.Error()
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpc.Do(req)
	if err != nil {
		return err.Error()
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err.Error()
	}
	return string(body)
}

func main() {
	base := env("PUZZLE_URL", "http://117.72.42.214:8080")
	httpc := &http.Client{Timeout: 30 * time.Second}

	// GET /api/ticket 领取通行证,本轮挑战全程复用同一张。
	resp, err := httpc.Get(base + "/api/ticket")
	if err != nil {
		fmt.Fprintln(os.Stderr, "领票失败:", err)
		os.Exit(1)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "读取通行证失败:", err)
		os.Exit(1)
	}
	jwt := firstJWT(string(body))
	if jwt == "" {
		fmt.Fprintln(os.Stderr, "未找到通行证")
		os.Exit(1)
	}

	// 并发提交编号 1..10;接收顺序既非完整升序也非完整降序即可通关。
	var newJWT string
	for round := 0; round < 5 && newJWT == ""; round++ {
		var wg sync.WaitGroup
		ch := make(chan string, 10)
		for seq := 1; seq <= 10; seq++ {
			wg.Add(1)
			go func(seq int) {
				defer wg.Done()
				ch <- postSeq(httpc, base, jwt, seq)
			}(seq)
		}
		wg.Wait()
		close(ch)
		for b := range ch {
			if strings.Contains(b, "保险箱已打开") {
				if j := firstJWT(b); j != "" {
					newJWT = j
				}
			}
		}
	}
	if newJWT == "" {
		fmt.Fprintln(os.Stderr, "未能通关")
		os.Exit(1)
	}

	// GET /api/reward 使用新通行证领取奖励。
	req, err := http.NewRequest(http.MethodGet, base+"/api/reward", nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "构造领奖请求失败:", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+newJWT)
	resp, err = httpc.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "领奖失败:", err)
		os.Exit(1)
	}
	body, err = io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, "读取奖励失败:", err)
		os.Exit(1)
	}
	fmt.Println(string(body))
}
