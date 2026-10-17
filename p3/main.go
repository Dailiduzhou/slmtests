package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Request struct {
	Op     string `json:"op"`
	Client string `json:"client"`
	Seq    int    `json:"seq"`
	Data   string `json:"data"`
}

// mailbox 表示一个信箱的状态与已存消息。
type mailbox struct {
	state    string // AWAKE / ASLEEP
	count    int
	nextSeq  int
	messages map[int]string
}

// Server 是小柠猫维护的消息服务器,按信箱标识管理各个信箱。
type Server struct {
	mailboxes map[string]*mailbox
}

func NewServer() *Server {
	return &Server{
		mailboxes: map[string]*mailbox{},
	}
}

// Handle 实现小柠猫对四种指令的响应,返回的 map 会自动被编码为 JSON。
func (s *Server) Handle(req Request) map[string]any {
	switch req.Op {
	case "MEOW":
		// 已存在信箱不能再次创建。
		if _, ok := s.mailboxes[req.Client]; ok {
			return map[string]any{"code": "ERR_ALREADY_EXISTS"}
		}
		s.mailboxes[req.Client] = &mailbox{state: "AWAKE", nextSeq: 1, messages: map[int]string{}}
		return map[string]any{"code": "PURR"}
	case "SEND":
		m, ok := s.mailboxes[req.Client]
		if !ok {
			return map[string]any{"code": "ERR_NO_SESSION"}
		}
		if m.state == "ASLEEP" {
			return map[string]any{"code": "ERR_ASLEEP"}
		}
		// 小柠猫喜欢把消息从 1 开始按顺序排序,前一编号消息存在时才能提交下一编号消息。
		if req.Seq == m.nextSeq {
			m.messages[req.Seq] = req.Data
			m.count++
			m.nextSeq++
			return map[string]any{"code": "ACK", "seq": req.Seq}
		}
		// 相同编号相同内容的消息成功提交返回同一成功响应;相同编号不同内容视为冲突。
		if stored, exists := m.messages[req.Seq]; exists {
			if stored == req.Data {
				return map[string]any{"code": "ACK", "seq": req.Seq}
			}
			return map[string]any{"code": "ERR_CONFLICT"}
		}
		return map[string]any{"code": "ERR_OUT_OF_ORDER"}
	case "PEEK":
		m, ok := s.mailboxes[req.Client]
		if !ok {
			return map[string]any{"code": "ERR_NO_SESSION"}
		}
		// 信箱关闭后仍可被查看。
		return map[string]any{"code": "PURR", "state": m.state, "count": m.count, "next_seq": m.nextSeq}
	case "SEAL":
		m, ok := s.mailboxes[req.Client]
		if !ok {
			return map[string]any{"code": "ERR_NO_SESSION"}
		}
		// 重复关闭信箱成功响应。
		m.state = "ASLEEP"
		return map[string]any{"code": "PURR"}
	default:
		return map[string]any{"code": "ERR_NOT_IMPLEMENTED"}
	}
}

func main() {
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	server := NewServer()

	for {
		var req Request
		if err := dec.Decode(&req); err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintln(os.Stderr, "读取请求失败:", err)
			return
		}
		if err := enc.Encode(server.Handle(req)); err != nil {
			fmt.Fprintln(os.Stderr, "输出响应失败:", err)
			return
		}
	}
}
