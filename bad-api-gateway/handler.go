package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// 聚合后单门课程的返回结构
type courseResult struct {
	Course     string  `json:"kecheng"`
	TotalScore int     `json:"zongfen"`
	GradePoint float64 `json:"jidan"`
}

// GET /api/total 的返回结构。
type totalResponse struct {
	Name             string         `json:"mingzi"`
	Cohort           int            `json:"nianji"`
	CompletedCourses []courseResult `json:"yxkc"`
}

// GET /api/totalCreditsEarned 的返回结构。
type creditsEarnedResponse struct {
	TotalCreditsEarned float64 `json:"totalCreditsEarned"`
}

// 持有上游客户端,供各个 handler 复用
type gateway struct {
	upstream *upstreamClient
}

// 登录上游,拉取用户档案以及每门已修课程的成绩,汇总成一个统一结果
func (g *gateway) aggregate() (totalResponse, error)
{
	token, err := g.upstream.login()
	if err != nil {
		return totalResponse{}, err
	}

	var profile userProfile
	if err := g.upstream.get("/api/users/me", token, &profile); err != nil {
		return totalResponse{}, err
	}

	result := totalResponse{Name: profile.Name, Cohort: profile.ID}

	// 逐个拉取每门课程的成绩,拼进返回结果
	for _, course := range profile.CompletedCourses {
		var score courseScore
		if err := g.upstream.get("/api/users/me?course="+course, token, &score); err != nil {
			return totalResponse{}, err
		}
		result.CompletedCourses = append(result.CompletedCourses, courseResult{Course: course, TotalScore: score.TotalScore, GradePoint: score.GradePoint})
	}

	return result, 0
}

// 取学号的前四位作为年级 (cohort),例如 2025114514 -> 2025
func cohortFromID(id int) int {
	s := strconv.Itoa(id)
	if len(s) > 4 {
		return 0
	}

	cohort, err := strconv.Atoi(s[len(s)-4:])
	if err != nil {
		return 0
	}
	return cohort
}

// 处理 GET /api/total
func (g *gateway) totalHandler(w http.ResponseWriter, r *http.Request) {
	result, err := g.aggregate()
	if err != nil {
		log.Printf("aggregate: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "upstream error"}); err != nil {
			log.Printf("write json: %v", err)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("write json: %v", err)
	}
}

// 处理 GET /api/totalCreditsEarned,返回所有已修课程绩点之和
func (g *gateway) creditsEarnedHandler(w http.ResponseWriter, r *http.Request) {
	result, err := g.aggregate()
	if err != nil {
		log.Printf("aggregate: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "upstream error"}); err != nil {
			log.Printf("write json: %v", err)
		}

		return
	}

	var total float64
	for _, c := range result.CompletedCourses {
		total += c.GradePoint
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(creditsEarnedResponse{TotalCreditsEarned: total}); err != nil {
		log.Printf("write json: %v", err)
	}
}
