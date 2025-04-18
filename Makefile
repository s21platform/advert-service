.PHONY: protogen

protogen:
	protoc --go_out=. --go-grpc_out=. ./api/advert.proto
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/advert.proto

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
