export GOPATH  = $(CURDIR)/../../lib/go

TEST_ARTIFACTS = \
  gen-go \
  test-compile-stamp \
  test-exercise-stamp \
  test-generation-stamp \
  test-stamp \
  test-validate-stamp

THRIFT=$(CURDIR)/../merge/thrift-0.8.0/compiler/cpp/thrift

all: test

test: test-stamp

test-stamp: test-exercise-stamp
	touch $@

test-exercise-stamp: test-compile-stamp
	cp -f simple_test.go gen-go/simple
	cd gen-go/simple && go test -v -x .
	touch $@

test-compile-stamp: test-validate-stamp
	cd gen-go/simple && go build -v -x ttypes.go
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
