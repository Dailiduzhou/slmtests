package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

// 预约状态取值。
const (
	statusBooking   = "booking"
	statusCancelled = "cancelled"
	statusUsed      = "used"
)

// Reservation 是一条 507 会议室预约记录。
type Reservation struct {
	ID      int    `json:"id"`
	Project string `json:"project"`
	Slot    string `json:"slot"`
	Status  string `json:"status"`
}

// booking 是 507 预约系统的内存存储,保护并发访问。
type booking struct {
	mu    sync.Mutex
	next  int
	items map[int]*Reservation
}

func newBooking() *booking {
	return &booking{next: 1, items: map[int]*Reservation{}}
}

// Reserve 创建预约,默认状态为 booking;同一时间段已有预约时失败。
func (b *booking) Reserve(project, slot string) (*Reservation, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.items {
		if r.Slot == slot && r.Status == statusBooking {
			return nil, false
		}
	}
	r := &Reservation{ID: b.next, Project: project, Slot: slot, Status: statusBooking}
	b.next++
	b.items[r.ID] = r
	return r, true
}

// Cancel 取消预约;已在使用中或已取消的预约不能取消。
func (b *booking) Cancel(id int) (*Reservation, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	r, ok := b.items[id]
	if !ok || r.Status != statusBooking {
		return nil, false
	}
	r.Status = statusCancelled
	return r, true
}

// Use 将预约标记为已使用(会议正在进行)。
func (b *booking) Use(id int) (*Reservation, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	r, ok := b.items[id]
	if !ok || r.Status != statusBooking {
		return nil, false
	}
	r.Status = statusUsed
	return r, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		fmt.Fprintln(os.Stderr, "写入响应失败:", err)
	}
}

// handleReservations 处理 POST /api/reservations 创建预约。
func (b *booking) handleReservations(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Project string `json:"project"`
		Slot    string `json:"slot"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "请求体不是有效的 JSON"})
		return
	}
	if body.Project == "" || body.Slot == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "预约需包含项目名与时间段"})
		return
	}
	r, ok := b.Reserve(body.Project, body.Slot)
	if !ok {
		writeJSON(w, http.StatusConflict, map[string]any{"error": "预约失败"})
		return
	}
	writeJSON(w, http.StatusCreated, r)
}

// handleReservationAction 处理 /api/reservations/{id}/{action} 的取消与使用。
func (b *booking) handleReservationAction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, action, ok := parseAction(req.URL.Path)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}
	var r *Reservation
	var done bool
	switch action {
	case "cancel":
		r, done = b.Cancel(id)
	case "use":
		r, done = b.Use(id)
	default:
		writeJSON(w, http.StatusNotFound, map[string]any{"error": "not found"})
		return
	}
	if !done {
		if action == "cancel" {
			writeJSON(w, http.StatusConflict, map[string]any{"error": "取消失败"})
		} else {
			writeJSON(w, http.StatusConflict, map[string]any{"error": "使用失败"})
		}
		return
	}
	writeJSON(w, http.StatusOK, r)
}

// parseAction 从 /api/reservations/{id}/{action} 解析编号与动作。
func parseAction(path string) (int, string, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 4 || parts[0] != "api" || parts[1] != "reservations" {
		return 0, "", false
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, "", false
	}
	return id, parts[3], true
}

// env 读取环境变量,未设置或为空时返回默认值。
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	store := newBooking()
	http.HandleFunc("/api/reservations", store.handleReservations)
	http.HandleFunc("/api/reservations/", store.handleReservationAction)
	addr := ":" + env("PORT", "8080")
	fmt.Fprintln(os.Stderr, "507 预约系统监听", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintln(os.Stderr, "启动服务失败:", err)
		os.Exit(1)
	}
}
