BINARY := todoist-proxy

.PHONY: build test run clean

build:
	go build -o $(BINARY) .

test:
	go test ./...

run: build
	TODOIST_PROXY_ALLOW=$(TODOIST_PROXY_ALLOW) PORT=$(PORT) ./$(BINARY)

clean:
	rm -f $(BINARY)
