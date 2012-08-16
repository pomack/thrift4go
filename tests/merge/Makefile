THRIFT_ARCHIVE_URL = https://dist.apache.org/repos/dist/release/thrift/0.8.0/thrift-0.8.0.tar.gz
TEST_ARTIFACTS = \
  test-compilation-stamp \
  test-installation-stamp \
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

test-stamp: test-installation-stamp
	touch $@

test-merge-stamp: thrift-extraction-stamp
	cd $(THRIFT4GO) && bash ./scripts/merge_and_build.sh -b
	touch $@

test-compilation-stamp: test-merge-stamp
	cd $(THRIFT) && ./configure \
		--with-go \
		--without-c_glib \
		--without-erlang \
		--without-haskell \
		--without-java \
		--without-perl \
		--without-php \
		--without-python \
		--without-ruby
	cd $(THRIFT) && make
	touch $@

test-installation-stamp: test-compilation-stamp tree-stamp
	cd $(THRIFT) && DESTDIR=$(CURDIR)/tree make install
	touch $@

tree-stamp:
	mkdir -vp tree
	touch $@

thrift-0.8.0.tar.gz:
	wget https://dist.apache.org/repos/dist/release/thrift/0.8.0/thrift-0.8.0.tar.gz

thrift-extraction-stamp: thrift-0.8.0.tar.gz
	tar xzvf thrift-0.8.0.tar.gz
	touch $@

clean:
	rm -rf $(TEST_ARTIFACTS)

.PHONY: test
