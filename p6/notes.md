# p6 解答笔记 — bad-api-gateway

## 修复内容

1. **编译错误**: `main.go` 引用了未定义的 `env`,新增 `env.go` 提供 `env(key, fallback)`(未设置或为空时返回默认值)。
2. **硬编码 → 环境变量**(应用启动时读取,未设置时用默认值):
   - `main.go`: 上游 URL 常量 `b` → `upstreamURL := env("UPSTREAM_URL", "https://tottering.muxixyz.com")`;端口已走 `env("PORT", "8080")`。
   - `client.go`: 测试账号常量 `upstreamUsername`/`upstreamPassword` → 包级变量,启动时从 `UPSTREAM_USER`/`UPSTREAM_PASS` 读取,默认 `wajueji`/`12345678`。
3. **逻辑错误**:
   - `client.go`: 认证头 `"Baerer"+token` 拼写错误 → `"Bearer "+token`。
   - `handler.go`: `Cohort` 直接取 `profile.ID` → 改为 `cohortFromID(profile.ID)`。
   - `cohortFromID`: 边界判断 `len(s) > 4` 方向错误 → `len(digits) < 4`;取末四位 `s[len(s)-4:]` → 取前四位 `digits[:4]`(2025114514 → 2025)。
   - `handler.go`: JSON 标签为拼音(`mingzi`/`nianji`/`yixiu`/`kecheng`/`zongfen`/`jidian`)→ 改为 `name`/`cohort`/`completedCourses`/`course`/`totalScore`/`gradePoint`,与 openapi.yaml 示例一致。
4. **复制粘贴 → 抽象新函数**: 两个 handler 重复的错误处理与 JSON 写出代码块,抽象为 `writeJSON(w, status, v)` 与 `writeUpstreamError(w, err)`。
5. **命名**: `b` → `upstreamURL`,`m` → `mux`,`g`(局部)→ `gw`,`s` → `digits`。
6. **格式化**: 重写 `client.go` 中挤成一行的 if 语句等,`gofmt -l .` 已无输出。

## 验证结果

- `go build ./...` 通过;`gofmt -l .` 无输出;`go vet ./...` 无警告(无死代码)。
- 上游 `https://tottering.muxixyz.com` 可达,实测:登录返回 token;`/api/users/me` 返回 name/id/completedCourses;`/api/users/me?course=math|c` 返回 60/1.5 与 100/2.5。
- 网关实测(`PORT=8081`):
  - `GET /api/total` → 200 `{"name":"wajueji","cohort":2025,"completedCourses":[{"course":"math","totalScore":60,"gradePoint":1.5},{"course":"c","totalScore":100,"gradePoint":2.5}]}`
  - `GET /api/totalCreditsEarned` → 200 `{"totalCreditsEarned":4}`
- 错误路径实测(`PORT=8092 UPSTREAM_URL=http://127.0.0.1:9`):两接口均返回 502 `{"error":"upstream error"}`。
- 默认端口:服务器成功监听 `:8080`(日志确认)。注意本机上 8080 已被用户自己的 Strata 服务占用,localhost:8080 的连接会被其截获;在无占用环境下默认端口即为题目所需的 Apifox 环境 URL。

## 收尾(终验)

- 按题目原文复核:`go vet .` 无警告、`gofmt -l .` 无输出、`go run .` 编译通过。
- 终验实测:`PORT=8081` 下两接口 200(JSON 与 openapi.yaml 示例逐字段一致);`PORT=8082 UPSTREAM_URL=http://127.0.0.1:9` 下两接口均 502 `{"error":"upstream error"}`。
- 仓库卫生:从提交中移除误入库的编译产物 `api-gateway`(约 10 MB 二进制),新增 `bad-api-gateway/.gitignore` 忽略该产物;amend 进 `feat: round 2` 后推送。
