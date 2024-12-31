SHELL := /bin/bash
.SHELLFLAGS := -euo pipefail -c ${SHELLFLAGS}
GINKGO_OPTS=-r --race --require-suite --randomize-all --cover ${OPTS}

# Metrics
cpu_app_name :=golang-autoscaler-cpuutil-metric
custom_metric_app_name :=golang-autoscaler-custom-metrics

.PHONY: build-for-custom-metrics
build: clean
	@echo "# building app ${custom_metric_app_name}"
	@go mod tidy
	# build for linux/amd64
	@env GOOS=linux GOARCH=amd64 go build -o deploy/build/${custom_metric_app_name}

.PHONY: build-for-cpu
build-for-cpu: clean
	@echo "# building app ${cpu_app_name}"
	@go mod tidy
	# build for linux/amd64
	@env GOOS=linux GOARCH=amd64 go build -o deploy/build/${cpu_app_name}

build-for-arm: clean
	@echo "# building app ${custom_metric_app_name}"
	@go mod tidy
	@go build -o deploy/build/${custom_metric_app_name}



# deploy app with cpuutil metric on cloud foundry
deploy-with-cpu: build-for-cpu
	@echo "# deploying app ${cpu_app_name}"
	@cf push -f deploy/app-manifest.yml -p deploy/build --no-start --var app_name=${cpu_app_name}
	@cf create-service autoscaler standard ak-test-autoscaler-${cpu_app_name}
	@cf bind-service ${cpu_app_name} ak-test-autoscaler-${cpu_app_name} -c deploy/cpu-utilization-policy.json
	@cf start ${cpu_app_name}

# remove all created resources in cf for app
clean-cpu-app:
	@echo "# removing all created resources in cf for app ${cpu_app_name}"
	@cf delete-service -f ak-test-autoscaler${cpu_app_name}
	@cf delete -f ${cpu_app_name}

test:
	@echo "Running tests"
	@ginkgo run ${GINKGO_OPTS}

buildtools:
	@echo "# Installing build tools"
	@go mod download
	@which ginkgo >/dev/null || go install github.com/onsi/ginkgo/v2/ginkgo


clean:
	@echo "# cleaning app"
	@rm -rf build
	@rm -rf deploy/build/* || true
	@mkdir -p deploy/build/
