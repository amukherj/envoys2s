GIT_REPO=envoys2s
CONTAINER_IMGS=rockers-img moods-img frontsvc-img edgesvc-img

all: bin/randlst bin/frontsvc $(CONTAINER_IMGS)

.PHONY: clean
clean:
	rm -rf bin/*
	docker rmi rockers-img:latest frontsvc-img:latest edgesvc-img:latest

bin/randlst:
	mkdir -p `dirname $@`
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $(GIT_REPO)/cmd/service

bin/frontsvc:
	mkdir -p `dirname $@`
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $(GIT_REPO)/cmd/`basename $@`

.PHONY: rockers-img
rockers-img: bin/randlst
	-docker rmi $@:latest >/dev/null 2>&1
	docker build -t $@:latest -f Dockerfile-rockers .

.PHONY: moods-img
moods-img:
	-docker rmi $@:latest >/dev/null 2>&1
	docker build -t $@:latest -f Dockerfile-moods .

.PHONY: frontsvc-img
frontsvc-img:
	-docker rmi $@:latest >/dev/null 2>&1
	docker build -t $@:latest -f Dockerfile-frontsvc .

.PHONY: edgesvc-img
edgesvc-img:
	-docker rmi $@:latest >/dev/null 2>&1
	docker build -t $@:latest -f Dockerfile-edgesvc .
