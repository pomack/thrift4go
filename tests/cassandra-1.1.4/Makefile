export GOPATH  = $(CURDIR)/../../lib/go

TEST_ARTIFACTS = \
	gen-go \
	test-compile-stamp \
	test-generation-stamp \
	test-stamp

THRIFT = $(CURDIR)/../merge/thrift-$(THRIFT_VERSION)/compiler/cpp/thrift
IDL = $(CURDIR)/cassandra.thrift

all: test

test: test-stamp

test-stamp: test-compile-stamp
	touch $@

test-compile-stamp: test-generation-stamp
	cd gen-go/cassandra && go build -v -x .
	touch $@

test-generation-stamp:
	"$(THRIFT)" --gen go "$(IDL)"
	touch $@

clean:
	rm -rf $(TEST_ARTIFACTS)

.PHONY: test
