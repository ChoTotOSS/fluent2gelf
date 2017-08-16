PACKAGE=github.com/ChoTotOSS/fluent2gelf
DIST_NAME=fluent2gelf

build:
	docker run -v `pwd`:/go/src/$(PACKAGE) golang:alpine go build -o /go/src/$(PACKAGE)/dist/$(DIST_NAME) -v $(PACKAGE)

docker: build
	docker build -t duythinht/fluent2gelf .
