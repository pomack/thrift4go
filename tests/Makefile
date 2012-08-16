SUB_DIRECTORIES = generation merge
CLEAN_DIRECTORIES = $(SUB_DIRECTORIES:%=%-clean)
TEST_STAMPS = $(SUB_DIRECTORIES:%=%-stamp)
TEST_ARTIFACTS =$(TEST_STAMPS)

all: test

test: $(SUB_DIRECTORIES)

generation: generation-stamp

generation-stamp: merge-stamp
	$(MAKE) -C generation test
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
