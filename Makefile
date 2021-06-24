build: ## 构建二进制
	#@bash hack/docker/build.sh ${version} ${tagversion} ${commit_sha1}
	# go get github.com/mitchellh/gox
	@gox -osarch="darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

lint:
	golangci-lint run  ./...


scan:
	go run main.go --key "drone-sonar-plugin"

default: lint scan