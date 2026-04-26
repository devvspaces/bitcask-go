client:
	go build -o ./bin/client ./cmd/main.go

kvstore:
	go build -o ./bin/kv ./kv/main.go

build:
	make client && make kvstore

runclient:
	go run ./cmd/main.go

runkv:
	go run ./kv/main.go