GOSRC = $(shell find . -type f -name '*.go')

VERSION=v1.0.0

build: controller

controller: $(GOSRC)
	CGO_ENABLED=0 GOOS=linux go build -o controller cmd/controller/controller.go


clean:
	rm -rf controller

clean-image:
	docker rmi linkingthing/web-controller:${VERSION}

.PHONY: clean install
