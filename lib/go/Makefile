export GOPATH = $(CURDIR)
ARTIFACTS = test-stamp

all: test

test: test-stamp

test-stamp:
	cd "$(GOPATH)" && go test -v -x thrift
	touch $@

clean:
	-rm -f $(ARTIFACTS)

.PHONY: test clean
