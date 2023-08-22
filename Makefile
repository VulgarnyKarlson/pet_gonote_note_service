GOCMD=go
GOBUILD=$(GOCMD) build
CMDPATH=cmd/main.go
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
PROJECTNAME=note_service
.PHONY: build build-linux debug race clean  run run-race run-docker

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o bin/main ${CMDPATH}
build:
	GO111MODULE=on $(GOBUILD) -o bin/main ${CMDPATH}
build-docker:
	docker build . -t $(PROJECTNAME) -f ./Dockerfile
debug:
	GO111MODULE=on $(GOBUILD) -gcflags="all=-N -l" -o bin/main ${CMDPATH}
race:
	GO111MODULE=on $(GOBUILD) -race -o bin/main ${CMDPATH}
clean:
	rm -f bin
run:
	make build && ./bin/main
run-race:
	GO111MODULE=on $(GOCMD) run -race ${CMDPATH}

run-docker:
	make build-docker \
	&&  docker stop $(PROJECTNAME);\
		docker rm $(PROJECTNAME);\
		docker run -d --network host --restart unless-stopped --name $(PROJECTNAME) -t $(PROJECTNAME) \
	&& docker logs -f $(PROJECTNAME)
