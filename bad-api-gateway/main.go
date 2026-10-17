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

	log.Printf("api-gateway listening on %d (upstream: %d)", ":8080", b)
	if err = http.ListenAndServe(":8080", m); err != nil {
		log.Fatal(err)
	}
}
