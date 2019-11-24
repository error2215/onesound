PROTO_ZIP='protoc-3.7.1-linux-x86_64.zip'

# Download and install all dependencies
install: mod
	@echo "All packages successfully installed!"
	@wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/$(PROTO_ZIP)
	@unzip -o $(PROTO_ZIP)
	@rm $(PROTO_ZIP)
	@ GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go
	@ GO111MODULE=on go build github.com/golang/protobuf/protoc-gen-go
	@ GO111MODULE=on go get github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	@ GO111MODULE=on go build github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

# Build gRPC files
build_go: clean_go
	@echo "Creating folders..."
	@echo "Generating common Go sources..."
	@./bin/protoc --plugin=protoc-gen-go -I. --go_out=compiled protobuf/common/*.proto
	@./bin/protoc --plugin=protoc-gen-go -I. --go_out=compiled protobuf/onesound_models/*.proto
	@echo "Generating client and server sources..."
	@./bin/protoc --plugin=protoc-gen-go -I. --go_out=plugins=grpc:compiled protobuf/onesound_api/*.proto  # also generate server class
	@echo "Moving generated files..."
	@mv -f compiled/onesound/compiled/* compiled/
	@rm -rf compiled/onesound


# Documentation
build_docs:
	@./bin/protoc --plugin=protoc-gen-doc=./protoc-gen-doc \
		--doc_out=./docs/ \
		--doc_opt=html,index.html \
		protobuf/*/*.proto
	@echo "Documentation generated! Don't forget to commit."

clean_go:
	rm -rf compiled/*

clean_docs:
	rm -rf docs/*

build: build_go build_docs build_bin

clean: clean_go clean_docs

mod:
	@echo "======================================================================"
	@echo "Run MOD"
	@ GO111MODULE=on go mod verify
	@ GO111MODULE=on go mod tidy
	@ GO111MODULE=on go mod vendor
	@ GO111MODULE=on go mod download
	@ GO111MODULE=on go mod verify
	@echo "======================================================================"

# Run tests
tests:
	$(SOURCE_PATH) go test  -coverprofile=coverage.out -v ./backend/api/
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

test: mod tests

# Build the server
build_bin:
	@echo "Building the binary..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./application ./main.go

# Run the server
run: build
	./application

