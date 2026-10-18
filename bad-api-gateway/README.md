# api-gateway

一个 API 网关, 从数据源网站 `https://tottering.muxixyz.com` 爬取数据,
再对下游暴露聚合后的接口, 用于方便地获取这些数据.

使用测试账号 `wajueji` / `12345678` 登录上游,
再拉取用户档案与各科成绩后汇总返回.

## 运行

```bash
go run .
# 默认监听 :8080
```

## 环境变量

- `PORT` — 监听端口,默认 `8080`
- `UPSTREAM_URL` — 上游地址,默认 `https://tottering.muxixyz.com`
- `UPSTREAM_USER` — 上游测试账号,默认 `wajueji`
- `UPSTREAM_PASS` — 上游测试密码,默认 `12345678`

## 接口

### `GET /api/total`

返回用户档案及每门已修课程的成绩,`cohort` 为学号前四位.

```json
{
  "name": "wajueji",
  "cohort": 2025,
  "completedCourses": [
    { "course": "math", "totalScore": 60, "gradePoint": 1.5 },
    { "course": "c", "totalScore": 100, "gradePoint": 2.5 }
  ]
}
```

### `GET /api/totalCreditsEarned`

返回所有已修课程绩点 (gradePoint) 之和.

```json
{ "totalCreditsEarned": 4 }
```

## 文件结构

| 文件 | 职责 |
| --- | --- |
| `main.go` | 程序入口、路由注册、监听端口 |
| `client.go` | 上游数据源客户端(登录、带认证 GET) |
| `handler.go` | 网关 handler、数据聚合、JSON 输出 |
| `env.go` | 读取环境变量,未设置时用默认值 |
