SHELL = /bin/bash

Branch=$(shell git symbolic-ref --short -q HEAD)
Commit=$(shell git rev-parse --short HEAD)
Date=$(shell git log --pretty=format:%cd $(Commit) -1)
Author=$(shell git log --pretty=format:%an $(Commit) -1)
shortDate=$(shell git log -1 --format="%at" | xargs -I{} date -d @{} +%Y%m%d)
Email=$(shell git log --pretty=format:%ae $(Commit) -1)
Ver=$(shell echo $(Branch)-$(Commit))
GoVersion=$(shell go version )
REPO ?= uchainorg/uscan

.PHONY: all-build
all-build: pull-submodule frontend-build statik build

.PHONY: pull-submodule
pull-submodule:
	git submodule update

.PHONY: frontend-build
frontend-build:
	cd uscan_frontend && yarn && yarn build

.PHONY: statik
statik:
	go install github.com/rakyll/statik@latest
	statik -f -src=uscan_frontend/dist

.PHONY: build
build: 
	go build -a -installsuffix cgo \
	-ldflags "-X 'github.com/uchainorg/uscan/cmd.Branch=$(Branch)' \
	-X 'github.com/uchainorg/uscan/cmd.Commit=$(Commit)' \
	-X 'github.com/uchainorg/uscan/cmd.Date=$(Date)' \
	-X 'github.com/uchainorg/uscan/cmd.Author=$(Author)' \
	-X 'github.com/uchainorg/uscan/cmd.Email=$(Email)' \
	-X 'github.com/uchainorg/uscan/cmd.GoVersion=$(GoVersion)'" -o bin/uscan

.PHONY: perf
perf: compile
	bin/uscan --config .uscan.yaml

.PHONY: race
race:
	go run -race main.go --config .uscan.yaml

.PHONY: start
start: compile
	bin/uscan --config .uscan.yaml


.PHONY: docker-build
docker-build:
	docker build \
	-t $(REPO):$(Ver) .
	docker tag $(REPO):$(Ver) $(REPO):latest
	docker image prune -f --filter label=stage=builder

.PHONY: docker-release
docker-release: build
	docker push $(REPO):$(Ver)
	docker push $(REPO):latest

