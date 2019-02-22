include Makefile.properties

start:
	@$(MAKE) deps
	@FORWARD_TO=http://localhost:1338 go run *.go
echoserver:
	@go run ./utils/echoserver/main.go
deps:
	@go mod vendor
compile:
	@$(MAKE) compile.linux
	@$(MAKE) compile.darwin
	@$(MAKE) compile.windows
compile.linux:
	@$(MAKE) BIN_NAME=$(BINARY_FILENAME) GOARCH=amd64 GOOS=linux .compile
compile.macos:
	@$(MAKE) BIN_NAME=$(BINARY_FILENAME) GOARCH=amd64 GOOS=darwin .compile
compile.windows:
	@$(MAKE) BIN_NAME=$(BINARY_FILENAME) GOARCH=386 GOOS=windows BIN_EXT=.exe .compile
.compile:
	@docker build \
		--target=development \
		--build-arg BIN_EXT=${BIN_EXT} \
		--build-arg BIN_NAME=${BIN_NAME} \
		--build-arg GOARCH=${GOARCH} \
		--build-arg GOOS=${GOOS} \
		-t $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		.
	-@docker stop proxy_gzip_for_binary_extraction && docker rm proxy_gzip_for_binary_extraction
	@docker run \
		-d \
		--entrypoint='sleep' \
		--name proxy_gzip_for_binary_extraction \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		1000
	@mkdir -p $(CURDIR)/bin
	@docker cp \
		proxy_gzip_for_binary_extraction:/go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} \
		$(CURDIR)/bin
	@docker cp \
		proxy_gzip_for_binary_extraction:/go/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}.sha256 \
		$(CURDIR)/bin
	@docker stop proxy_gzip_for_binary_extraction && docker rm proxy_gzip_for_binary_extraction
	@rm -rf $(CURDIR)/bin/$(BINARY_FILENAME)
	@chmod +x $(CURDIR)/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT}
	@ln -s $(CURDIR)/bin/${BIN_NAME}-${GOOS}-${GOARCH}${BIN_EXT} $(CURDIR)/bin/$(BINARY_FILENAME)
	@chmod +x $(CURDIR)/bin/$(BINARY_FILENAME)
package:
	@$(MAKE) package.docker
package.docker:
	@docker build \
		--target=production \
		-t $(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest \
		.
release:
	@$(MAKE) release.docker
release.docker:
	@$(MAKE) version.get > .version
	@docker push \
		$(DOCKER_REGISTRY_HOSTNAME)/$(DOCKER_IMAGE_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest
version.get:
	@docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		get-latest -q
version.bump:
	@docker run \
		-v "$(CURDIR):/app" \
		zephinzer/vtscripts:latest \
		iterate -i -q
