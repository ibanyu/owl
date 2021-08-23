test: fmt
	go test -race  ./controller/test/...
	go test -race  ./service/checker/...
	go test -race  ./service/sql_util/...
	go test -race  ./util/...

build: fmt
	mkdir -p bin
	go build -o bin/owl ./cmd/owl/

build-linux: fmt
	CGO_ENABLED=0 GOOS=linux go build -o bin/owl -a -ldflags '-extldflags "-static"' ./cmd/owl/

fmt:
	go fmt ./...

run: build
	./bin/owl

.ONESHELL:
build-front:
	mkdir -p bin/front
	if [ ! -e "./bin/front/owl_web" ]; then cd bin/front && git clone https://github.com/ibanyu/owl_web.git; else cd bin/front/owl_web && git pull; fi
	cd bin/front/owl_web && yarn install && yarn build
	rm -rf ./static && mkdir static && mv bin/front/owl_web/dist/* ./static/

build-docker: build
	docker build -t palfish/owl:v0.1.0 .

run-docker: build-docker
	docker run -p 8081:8081 -d  palfish/owl:v0.1.0
