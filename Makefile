.PHONY: gen-demo-proto
gen-demo-proto:
	@cd demo/demo_proto && cwgo server -I ../../idl --module github.com/cloudwego/biz-demo/gomall/demo/demo_proto --service demo_proto --idl ../../idl/echo.proto

.PHONY: gen-demo-thrift
gen-demo-thrift:
	@cd demo/demo_thrift && cwgo server --module github.com/cloudwego/biz-demo/gomall/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift

.PHONY: demo-link-fix
demo-link-fix:
	cd demo/demo_proto && golangci-lint run -E gofumpt --path-prefix=. --fix --timeout=5m

.PHONY: gen-frontend
gen-frontend:
	@cd app/frontend && cwgo server -I ../../idl --type HTTP --service frontend --module gomall/app/frontend --idl ../../idl/frontend/home.proto

.PHONY: gen-user
gen-user:
	@cd app/user && cwgo server --type RPC --service user --module gomall/app/user --pass "-use gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/user.proto
	@cd rpc_gen && cwgo client -I ../idl --type RPC --service user --module gomall/rpc_gen --idl ../idl/user.proto

