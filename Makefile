GOCMD=go
GOBUILD=$(GOCMD) build
CMDNOTEPATH=cmd/service_note/main.go
CMDPRODUCERPATH=cmd/service_stats_sender/main.go
BINNOTEPATH=bin/service_note
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
PROJECTNAME=note_service
.PHONY: build build-linux debug race clean  run run-race run-docker

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ${BINNOTEPATH} ${CMDNOTEPATH} && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o bin/service_stats_sender ${CMDPRODUCERPATH}
build:
	make gen && \
	GO111MODULE=on $(GOBUILD) -o ${BINNOTEPATH} ${CMDNOTEPATH} && \
	GO111MODULE=on $(GOBUILD) -o bin/service_stats_sender ${CMDPRODUCERPATH}
build-docker:
	docker build . -t $(PROJECTNAME) -f ./DockerfileNote
debug:
	GO111MODULE=on $(GOBUILD) -gcflags="all=-N -l" -o ${BINNOTEPATH} ${CMDNOTEPATH}
race:
	GO111MODULE=on $(GOBUILD) -race -o ${BINNOTEPATH} ${CMDNOTEPATH}
clean:
	rm -f bin
gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/auth.proto
run:
	make build && ./${BINNOTEPATH}
run-race:
	GO111MODULE=on $(GOCMD) run -race ${CMDNOTEPATH}

run-docker-note:
	make build-docker \
	&&  docker stop $(PROJECTNAME);\
		docker rm $(PROJECTNAME);\
		docker run -d --network host --restart unless-stopped --name $(PROJECTNAME) -t $(PROJECTNAME) \
	&& docker logs -f $(PROJECTNAME)

run-docker-producer:
	make build-docker \
	&&  docker stop $(PROJECTNAME);\
		docker rm $(PROJECTNAME);\
		docker run -d --network host --restart unless-stopped --name $(PROJECTNAME) -t $(PROJECTNAME) \
	&& docker logs -f $(PROJECTNAME)
