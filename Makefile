.PHONY: build
build: 
	go build -o output/sqlgogen main.go

.PHONY: test
test:
	go get -u github.com/rakyll/gotest; \
	gotest -v ./... -count=1;