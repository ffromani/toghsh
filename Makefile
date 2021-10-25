all: toghsh

.PHONY: clan
clean:
	rm -rf _out

.PHONY: clean-deps
clean-deps:
	rm -rf vendor

.PHONY: update-deps
update-deps:
	go mod tidy && go mod vendor

toghsh-static: outdir
	CGO_ENABLED=0 go build -o _out/toghsh ./cmd/toghsh

toghsh: outdir
	go build -o _out/toghsh ./cmd/toghsh/

outdir:
	@mkdir -p _out || :

.PHONY: test-unit
test-unit:
	go test ./pkg/...

.PHONY: gofmt
gofmt:
	@echo "Running gofmt"
	gofmt -s -w `find . -path ./vendor -prune -o -type f -name '*.go' -print`
