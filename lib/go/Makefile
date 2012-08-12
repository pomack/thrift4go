export GOPATH = $(CURDIR)
ARTIFACTS = test-stamp

all: test

test: test-stamp

test-stamp:
	cd "$(GOPATH)" && go build -v -x thrift
	cd "$(GOPATH)" && go test -v -x thrift
	touch $@

clean:
	-rm -f $(ARTIFACTS)

.PHONY: test clean
