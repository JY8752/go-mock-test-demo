BIN := .tmp/bin

generate: | $(BIN)/moq
	@PATH=$(abspath $(@D))/$(BIN):$(PATH); go generate ./...

$(BIN):
	@mkdir -p $(BIN)

$(BIN)/moq: $(BIN) Makefile
	GOBIN=$(abspath $(@D)) go install github.com/matryer/moq@latest