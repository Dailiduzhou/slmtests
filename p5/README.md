# p5 「Muxi507-booking」507 预约系统

muxi 有很多项目小组,当某个项目需要线下开会时,他们会前往 507。
为避免多个小组同时到达 507 造成占用冲突,wjj 学长设计了 507 预约系统。

## 运行

```bash
go run .
# 默认监听 :8080,可用环境变量 PORT 覆盖
PORT=18050 go run .
```

## 接口

### 预约 `POST /api/reservations`

请求体需包含项目名与时间段,预约状态默认为 `booking`:

```json
{"project": "蓝桥杯", "slot": "2026-10-17 14:00"}
```

- 成功:`201` 返回预约的全部信息
  ```json
  {"id": 1, "project": "蓝桥杯", "slot": "2026-10-17 14:00", "status": "booking"}
  ```
- 同一时间段已有预约:`409` `{"error": "预约失败"}`

### 取消预约 `POST /api/reservations/{id}/cancel`

- 成功:`200` 预约状态修改为 `cancelled`
- 预约已在使用中(或已取消):`409` `{"error": "取消失败"}`

### 标记使用 `POST /api/reservations/{id}/use`

将预约状态修改为 `used`(会议进行中),便于演示"已在使用中的预约不能取消"。

## 状态模型

`booking`(已预约) → `cancelled`(已取消) 或 `used`(已使用)。
只有 `booking` 状态的预约可以被取消或使用。

## 测试

```bash
bash test.sh
```
