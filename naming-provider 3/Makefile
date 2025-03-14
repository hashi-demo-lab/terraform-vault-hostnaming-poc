default: install

.PHONY: build
build:
	go build -o terraform-provider-naming

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/local/terraform-providers/naming/1.0.0/darwin_arm64
	cp terraform-provider-naming ~/.terraform.d/plugins/local/terraform-providers/naming/1.0.0/darwin_arm64/

.PHONY: clean
clean:
	rm -f terraform-provider-naming

.PHONY: test
test:
	go test -v ./...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v ./... 