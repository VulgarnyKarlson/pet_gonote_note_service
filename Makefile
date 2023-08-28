GOCMD=go
GOBUILD=GO111MODULE=on $(GOCMD) build
CMDNOTEPATH=cmd/service_note/main.go
CMDPRODUCERPATH=cmd/service_stats_sender/main.go
BINNOTEPATH=bin/service_note
BINPRODUCERPATH=bin/service_stats_sender
CGO_ENABLED=0
GOOS=linux
GOARCH=amd64
PROJECTNAME=note_service
.PHONY: build build-linux debug race clean  run run-race run-docker

build-note-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ${BINNOTEPATH} ${CMDNOTEPATH}
build-producer-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-s -w" -o ${BINPRODUCERPATH} ${CMDPRODUCERPATH}
build-note:
	make gen && \
	$(GOBUILD) -o ${BINNOTEPATH} ${CMDNOTEPATH}
build-producer:
	make gen && \
	$(GOBUILD) -o ${BINPRODUCERPATH} ${CMDPRODUCERPATH}
build-docker-note:
	docker build . -t $(PROJECTNAME) -f ./DockerfileNote
build-docker-producer:
	docker build . -t $(PROJECTNAME) -f ./DockerfileProducer
race:
	GO111MODULE=on $(GOBUILD) -race -o ${BINNOTEPATH} ${CMDNOTEPATH}
clean:
	rm -f bin
gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/auth.proto
run-note:
	make build-note && ./${BINNOTEPATH}
run-producer:
	make build-producer && ./${BINPRODUCERPATH}
run-note-race:
	$(GOCMD) run -race ${CMDNOTEPATH}
run-producer-race:
	$(GOCMD) run -race ${CMDPRODUCERPATH}

run-docker-note:
	make build-docker-note \
	&&  docker stop $(PROJECTNAME);\
		docker rm $(PROJECTNAME);\
		docker run --network host --restart unless-stopped --name $(PROJECTNAME) -t $(PROJECTNAME) \
	&& docker logs -f $(PROJECTNAME)

run-docker-producer:
	make build-docker-producer \
	&&  docker stop $(PROJECTNAME);\
		docker rm $(PROJECTNAME);\
		docker run --network host --restart unless-stopped --name $(PROJECTNAME) -t $(PROJECTNAME) \
	&& docker logs -f $(PROJECTNAME)
