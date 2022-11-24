SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c ${SHELLFLAGS}
GINKGO_OPTS=-r --race --require-suite --randomize-all --cover ${OPTS}


.PHONY: build
build:
	echo "# building app"
	rm -rf deploy/build/* || true
	mkdir -p deploy/build/
	go mod tidy
	# build for linux/amd64
	env GOOS=linux GOARCH=amd64 go build -o deploy/build/golang-autoscaler-custom-metrics

test:
	@echo "Running tests"
	ginkgo run ${GINKGO_OPTS}

buildtools:
	@echo "# Installing build tools"
	@go mod download
	@which ginkgo >/dev/null || go install github.com/onsi/ginkgo/v2/ginkgo


clean:
	@echo "# cleaning app"
	@rm -rf build