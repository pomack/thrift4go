THRIFT_ARCHIVE_URL = https://dist.apache.org/repos/dist/release/thrift/0.8.0/thrift-0.8.0.tar.gz
TEST_ARTIFACTS = \
  test-merge-stamp \
  test-stamp \
  thrift-0.8.0 \
  thrift-0.8.0.tar.gz \
  thrift-extraction-stamp \
  tree \
  tree-stamp

export THRIFT=$(CURDIR)/thrift-0.8.0
export THRIFT4GO=$(CURDIR)/../..

all: test

test: test-stamp

test-stamp: test-merge
	touch $@

test-merge: test-merge-stamp

test-merge-stamp: thrift-extraction tree-stamp
	cd $(THRIFT4GO) && bash ./scripts/merge_and_build.sh -b
	cd $(THRIFT) && ./configure \
		--prefix=$(CURDIR)/tree \
		--with-go \
		--without-{c_glib,java,haskell,php,erlang,perl,python,ruby}
	cd $(THRIFT) && make
	touch $@

tree-stamp:
	mkdir -vp tree
	touch $@

thrift-extraction: thrift-extraction-stamp

thrift-0.8.0.tar.gz:
	wget https://dist.apache.org/repos/dist/release/thrift/0.8.0/thrift-0.8.0.tar.gz

thrift-extraction-stamp: thrift-0.8.0.tar.gz
	tar xzvf thrift-0.8.0.tar.gz
	touch $@

clean:
	rm -rf $(TEST_ARTIFACTS)

.PHONY: test test-merge thrift-extraction
