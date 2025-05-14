genAuth:
	@protoc \
	--go_out=proto \
	--go-grpc_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	auth\auth.proto

genOrder:
	@protoc \
	--go_out=proto \
	--go-grpc_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	order\order.proto

genCourier:
	@protoc \
	--go_out=proto \
	--go-grpc_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	courier\courier.proto

genDelivery:
	@protoc \
	--go_out=proto \
	--go-grpc_out=proto \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	delivery\delivery.proto \
	order\order.proto \
	courier\courier.proto
