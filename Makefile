gen: ## Generate Code
	goctl rpc protoc wallet.proto --go_out=.  --go-grpc_out=.  --zrpc_out=.