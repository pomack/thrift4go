export GOPATH  = $(CURDIR)/../../lib/go

TEST_ARTIFACTS = \
  gen-go \
  test-compile-stamp \
  test-generation-stamp \
  test-stamp \
  test-validate-stamp

THRIFT=$(CURDIR)/../merge/thrift-0.8.0/compiler/cpp/thrift

all: test

test: test-stamp

test-stamp: test-compile-stamp
	touch $@

test-compile-stamp: test-validate-stamp
	cd gen-go/simple && go build ttypes.go
	touch $@

test-validate-stamp: test-generation-stamp
	diff -Nru whitebox/gen-go/ gen-go/
	touch $@

test-generation-stamp:
	$(THRIFT) --gen go simple.thrift
	touch $@

clean:
	rm -rf $(TEST_ARTIFACTS)

.PHONY: test
