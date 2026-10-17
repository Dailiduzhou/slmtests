package main

import (
	"log"
	"net/http"
)

// 数据源网站的默认地址
const b = "https://tottering.muxixyz.com"

func main() {
	g := &gateway{upstream: newUpstreamClient(b)}

	// net/http 路由
	m := http.NewServeMux()
	m.HandleFunc("/api/total", g.totalHandler)
	m.HandleFunc("/api/totalCreditsEarned", g.creditsEarnedHandler)

	addr := ":" + env("PORT", "8080")
	log.Printf("api-gateway listening on %s (upstream: %s)", addr, b)
	if err := http.ListenAndServe(addr, m); err != nil {
		log.Fatal(err)
	}
}
