# p4 「🍥找回计划」解题记录

## 探索路径

1. `GET http://117.72.42.214:8080/`
   - 正文只有幌子,响应头有 `X-Next-Step: /api/ticket`。
   - 提示:网络响应不只有正文,要看"头"部。
2. `GET /api/ticket`
   - 返回 JWT 通行证(`Authorization: Bearer <JWT>`),有效期 30 分钟。
   - 每次领票生成独立会话,只领一次,本轮复用。
   - 响应头 `X-Next-Step: /api/challenge`。
3. `POST /api/challenge`,请求体 `{"seq": N}`,携带 `Authorization: Bearer <通行证>` 与 `Content-Type: application/json`。
   - 需要提交编号 1~10,且服务端接收顺序**既不是完整升序,也不是完整降序**。
   - 本次采用顺序 `1, 3, 5, 7, 9, 2, 4, 6, 8, 10`(交错序)通关。
   - 通关响应给出新通行证(有效期 10 分钟)与奖励入口 `/api/reward`。
4. `GET /api/reward`,携带新通行证 → 展示 🍥 与密钥。

## 答案

- 🍥:已找回 wjj 学长的鱼饼,挑战成功。
- 密钥:`d2pq5a2m6ZW/54ix6bG86aW855qE`
- 密钥为 Base64 编码,解码结果为:`wjj学长爱鱼饼的`

```bash
echo 'd2pq5a2m6ZW/54ix6bG86aW855qE' | base64 -d
# wjj学长爱鱼饼的
```

## 求解器

`main.go` 用 Go 实现:领票一次 → `go` 关键字并发提交编号 1~10 → 等待全部响应并读取 → 通关后用新通行证领奖。
环境变量 `PUZZLE_URL` 可覆盖站点地址,默认 `http://117.72.42.214:8080`。
