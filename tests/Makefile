SUB_DIRECTORIES = cassandra-1.0.11 cassandra-1.1.4 generation interoperability merge
CLEAN_DIRECTORIES = $(SUB_DIRECTORIES:%=%-clean)
TEST_STAMPS = $(SUB_DIRECTORIES:%=%-stamp)
TEST_ARTIFACTS =$(TEST_STAMPS)

all: test

test: $(SUB_DIRECTORIES)

cassandra-1.0.11: cassandra-1.0.11-stamp

cassandra-1.0.11-stamp: interoperability-stamp
	$(MAKE) -C cassandra-1.0.11 test
	touch $@

cassandra-1.1.4: cassandra-1.1.4-stamp

cassandra-1.1.4-stamp: interoperability-stamp
	$(MAKE) -C cassandra-1.1.4 test
	touch $@

generation: generation-stamp

generation-stamp: merge-stamp
	$(MAKE) -C generation test
	touch $@

interoperability: interoperability-stamp

interoperability-stamp: generation-stamp
	$(MAKE) -C interoperability test
	touch $@

merge: merge-stamp

merge-stamp:
	$(MAKE) -C merge test
	touch $@

clean: $(CLEAN_DIRECTORIES)
	rm -f $(TEST_ARTIFACTS)

$(CLEAN_DIRECTORIES):
	$(MAKE) -C $(@:%-clean=%) clean

.PHONY: test $(SUB_DIRECTORIES) $(CLEAN_DIRECTORIES)
