export GOPATH  = $(CURDIR)/../../lib/go

TEST_ARTIFACTS = \
	gen-go \
	gen-java \
	src/main/java/thrift4go/generated/* \
	target \
	test-compile-stamp \
	test-exercise-stamp \
	test-generation-stamp \
	test-stamp \
	test-start-binary-server-stamp \
	test-start-compact-server-stamp \
	test-start-json-server-stamp

THRIFT = $(CURDIR)/../merge/thrift-0.8.0/compiler/cpp/thrift
IDL = $(CURDIR)/../generation/simple.thrift
GENERATED = gen-java/thrift4go/generated/*.java
GENERATED_TARGET = src/main/java/thrift4go/generated

all: test

test: test-stamp

test-stamp: test-exercise-stamp
	touch $@

test-exercise-stamp: test-start-binary-server-stamp test-start-json-server-stamp test-start-compact-server-stamp exercise_test.go
	cp exercise_test.go gen-go/simple
	cd gen-go/simple && go test -v -x .
	touch $@

test-start-binary-server-stamp: test-compile-stamp
	java -jar target/interoperability-unversioned-jar-with-dependencies.jar BINARY 8080 &
	sleep 5
	touch $@

test-start-json-server-stamp: test-compile-stamp
	java -jar target/interoperability-unversioned-jar-with-dependencies.jar JSON 8081 &
	sleep 5
	touch $@

test-start-compact-server-stamp: test-compile-stamp
	java -jar target/interoperability-unversioned-jar-with-dependencies.jar COMPACT 8082 &
	sleep 5
	touch $@

test-compile-stamp: test-generation-stamp
	mvn compile
	mvn assembly:assembly
	mvn assembly:single
	touch $@

test-generation-stamp:
	"$(THRIFT)" --gen java "$(IDL)"
	"$(THRIFT)" --gen go "$(IDL)"
	mkdir -vp "$(GENERATED_TARGET)"
	mv $(GENERATED) "$(GENERATED_TARGET)"
	touch $@

clean:
	rm -rf $(TEST_ARTIFACTS)
	-pkill -f 8080
	-pkill -f 8081
	-pkill -f 8082
	-pkill -f 8083

.PHONY: test
