export ROOT_MOD=gomall
.PHONY: gen-demo-proto
gen-demo-proto:
	@cd demo/demo_proto && cwgo server -I ../../idl --module ${ROOT_MOD}/demo/demo_proto --service demo_proto --idl ../../idl/echo.proto

.PHONY: gen-demo-thrift
gen-demo-thrift:
	@cd demo/demo_thrift && cwgo server --module ${ROOT_MOD}/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift

.PHONY: demo-link-fix
demo-link-fix:
	cd demo/demo_proto && golangci-lint run -E gofumpt --path-prefix=. --fix --timeout=5m

.PHONY: gen-frontend
gen-frontend:
	@cd app/frontend && cwgo server -I ../../idl --type HTTP --service frontend --module ${ROOT_MOD}/app/frontend --idl ../../idl/frontend/home.proto

.PHONY: gen-user
gen-user:
	@cd app/user && cwgo server --type RPC --service user --module ${ROOT_MOD}/app/user --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/user.proto
	@cd rpc_gen && cwgo client -I ../idl --type RPC --service user --module ${ROOT_MOD}/rpc_gen --idl ../idl/user.proto


.PHONY: gen-product
gen-product:
	@cd app/product && cwgo server --type RPC --service product --module ${ROOT_MOD}/app/product --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/product.proto
	@cd rpc_gen && cwgo client -I ../idl --type RPC --service product --module ${ROOT_MOD}/rpc_gen --idl ../idl/product.proto


.PHONY: gen-seckill
gen-seckill:
	@cd app/seckill && cwgo server --type RPC --service seckill --module ${ROOT_MOD}/app/seckill --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/seckill.proto
	@cd rpc_gen && cwgo client -I ../idl --type RPC --service seckill --module ${ROOT_MOD}/rpc_gen --idl ../idl/seckill.proto
