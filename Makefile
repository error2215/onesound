PROTO_ZIP='protoc-3.7.1-linux-x86_64.zip'

# Download and install all dependencies
install:
	@wget https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/$(PROTO_ZIP)
	@unzip -o $(PROTO_ZIP)
	@rm $(PROTO_ZIP)
	@go get github.com/golang/protobuf/protoc-gen-go
	@go build github.com/golang/protobuf/protoc-gen-go
	@go get github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	@go build github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

# Build gRPC files
build_go: clean_go
	@echo "Creating folders..."
	@mkdir compiled
	@mkdir docs
	@echo "Generating common Go sources..."
	@./bin/protoc --plugin=protoc-gen-go -I. --go_out=compiled protobuf/common/*.proto
	@echo "Generating client and server sources..."
	@./bin/protoc --plugin=protoc-gen-go -I. --go_out=plugins=grpc:compiled protobuf/onesound_api/*.proto  # also generate server class
	@echo "Moving generated files..."
	@mv -f compiled/onesound/protobuf/* compiled/
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

build: build_go build_docs

clean: clean_go clean_docs