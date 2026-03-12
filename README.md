## wallet-srv

一个简单的 gRPC 钱包服务示例（创建钱包、查询余额、转账）。当前实现使用**进程内内存存储**，重启会丢数据，仅用于演示/练习。

## 前置依赖

- **Go**：请安装本机 Go（建议 1.22+）。安装后应能执行 `go version`
- （可选）**grpcurl**：用于本地快速调用 gRPC 接口验证
- （可选）**goctl**：如果你需要重新从 `wallet.proto` 生成代码

## 配置

默认配置文件：`etc/walletsrv.yaml`

```yaml
Name: walletsrv.rpc
Mode: dev
ListenOn: 0.0.0.0:8080
```

你也可以通过 `-f` 指定配置文件路径。

## 启动服务

在项目根目录执行：

```bash
go run ./cmd/walletsrv.go -f etc/walletsrv.yaml
```

看到类似输出说明启动成功：

```text
Starting rpc server at 0.0.0.0:8080...
```

## 本地调用验证（grpcurl）

开发/测试模式会启用 gRPC reflection（见 `cmd/walletsrv.go`），因此可以直接用 `grpcurl` 调用。

### 1) 查看服务与方法

```bash
grpcurl -plaintext 127.0.0.1:8080 list
grpcurl -plaintext 127.0.0.1:8080 list wallet_srv.Wallet
```


### 2) 创建钱包（Wallets）

```bash
grpcurl -plaintext -d '{"UserId":1}' 127.0.0.1:8080 wallet_srv.Wallet/Wallets
grpcurl -plaintext -d '{"UserId":2}' 127.0.0.1:8080 wallet_srv.Wallet/Wallets
```

### 3) 查询钱包（WalletsByID）

```bash
grpcurl -plaintext -d '{"walletID":1}' 127.0.0.1:8080 wallet_srv.Wallet/WalletsByID
```

### 4) 转账（Transfer）

`amount` 为字符串（decimal）。

```bash
grpcurl -plaintext -d '{"fromWalletID":1,"toWalletID":2,"amount":"10"}' 127.0.0.1:8080 wallet_srv.Wallet/Transfer
```

再次查询余额：

```bash
grpcurl -plaintext -d '{"walletID":1}' 127.0.0.1:8080 wallet_srv.Wallet/WalletsByID
grpcurl -plaintext -d '{"walletID":2}' 127.0.0.1:8080 wallet_srv.Wallet/WalletsByID
```

> 注意：当前实现没有“充值/入金”接口，默认余额为 0；因此去掉了余额的验证。

## 运行测试

运行全部单测：

```bash
go test ./... -v
```

如果你想看竞态（建议本机跑）：

```bash
go test ./... -race
```

## 代码生成（可选）

仓库包含 `Makefile` 里的生成命令：

```bash
make gen
```

它等价于（需要 `goctl` 与 `protoc` 环境）：

```bash
goctl rpc protoc wallet.proto --go_out=. --go-grpc_out=. --zrpc_out=.
```

## 目录结构（简要）

- `cmd/walletsrv.go`：服务入口（读取配置、启动 gRPC）
- `internal/server/`：gRPC server（由 goctl 生成）
- `internal/logic/`：业务逻辑层
- `internal/dao/`：存储/数据访问层（当前为内存实现）
- `etc/`：配置文件
- `wallet_srv/`：proto 生成的 gRPC 代码
