default: fmt lint install generate

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

generate:
	cd tools; go generate ./...

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

# Generate documentation
.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Run tests
.PHONY: test
test:
	go test -v ./...

# Build provider
.PHONY: build
build:
	go build -o terraform-provider-centreon

# Install provider locally
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/hashicorp/centreon/0.1.0/darwin_amd64/
	mv terraform-provider-centreon ~/.terraform.d/plugins/registry.terraform.io/hashicorp/centreon/0.1.0/darwin_amd64/

# Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-centreon

.PHONY: fmt lint test testacc build install generate docs clean
