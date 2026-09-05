default: fmt lint build generate

build:
	go build -v ./...

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

# Generate documentation
.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs --provider-name centreon --rendered-provider-name Centreon

.PHONY: fmt lint test testacc build generate docs
