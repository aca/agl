PKGS = $(shell go list ./... | grep -v /vendor/ | grep -v /bindata)

cover:
	@go test -coverprofile=cover.out -coverpkg=./... ./...
	@go tool cover -html=cover.out

wasm:
	GOOS=js GOARCH=wasm go build -o docs/main.wasm main.go

count:
	@find \
		./pkg/agl \
		./scripts \
		-name '*.go' \
		| xargs wc -l \
		| sort