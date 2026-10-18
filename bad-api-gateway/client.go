package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 上游数据源网站测试账号,启动时从环境变量读取,未设置时用默认值。
var (
	upstreamUsername = env("UPSTREAM_USER", "wajueji")
	upstreamPassword = env("UPSTREAM_PASS", "12345678")
)

// 上游 /api/users/me 返回的用户档案
type userProfile struct {
	Name             string   `json:"name"`
	ID               int      `json:"id"`
	CompletedCourses []string `json:"completedCourses"`
}

// 上游 /api/users/me?course=xxx 返回的单门课程成绩
type courseScore struct {
	TotalScore int     `json:"totalScore"`
	GradePoint float64 `json:"gradePoint"`
}

// 负责调用数据源网站的 Client
type upstreamClient struct {
	baseURL string
	http    *http.Client
}

// 创建一个指向 baseURL 的上游客户端
func newUpstreamClient(baseURL string) *upstreamClient {
	return &upstreamClient{
		baseURL: baseURL,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// 登录上游,成功后返回 Bearer token
func (c *upstreamClient) login() (string, error) {
	body, err := json.Marshal(map[string]string{"username": upstreamUsername, "password": upstreamPassword})
	if err != nil {
		return "", err
	}
	resp, err := c.http.Post(c.baseURL+"/api/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("调用登录接口失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("登录失败: 上游返回状态码 %d", resp.StatusCode)
	}
	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析登录响应失败: %w", err)
	}
	return result.Token, nil
}

// 发送一个带认证的 GET 请求,并把响应的 JSON 解析到 out 里
func (c *upstreamClient) get(path, token string, out any) error {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("请求 %s 失败: %w", path, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求 %s 失败: 上游返回状态码 %d", path, resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("解析 %s 响应失败: %w", path, err)
	}
	return nil
}
