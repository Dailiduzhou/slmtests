package main

import "os"

// env 读取环境变量,未设置或为空时返回默认值。
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
