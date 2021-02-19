LDFLAGS := "-s -w"

.PHONY: all
all: local

.PHONY: local
local:
	go build

.PHONY: dist
dist:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/ctrn

