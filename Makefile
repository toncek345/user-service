.PHONY: gen_proto
gen_proto:
	@protoc  -I=. \
		-I=./googleapis \
		--proto_path=./proto \
		--go-grpc_out=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--grpc-gateway_out . \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		--grpc-gateway_opt generate_unbound_methods=true \
		--openapiv2_out . \
		--openapiv2_opt logtostderr=true \
		`find ./proto -name '*.proto'`
