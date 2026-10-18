package main

import (
	"log"
	"net/http"
)

func main() {
	// 数据源网站的默认地址,启动时从环境变量读取
	upstreamURL := env("UPSTREAM_URL", "https://tottering.muxixyz.com")
	gw := &gateway{upstream: newUpstreamClient(upstreamURL)}

	// net/http 路由
	mux := http.NewServeMux()
	mux.HandleFunc("/api/total", gw.totalHandler)
	mux.HandleFunc("/api/totalCreditsEarned", gw.creditsEarnedHandler)

	addr := ":" + env("PORT", "8080")
	log.Printf("api-gateway listening on %s (upstream: %s)", addr, upstreamURL)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
