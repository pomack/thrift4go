include Makefile.environment

SUB_DIRECTORIES = lib tests
CLEAN_DIRECTORIES = $(SUB_DIRECTORIES:%=clean-%)
TEST_DIRECTORIES = $(SUB_DIRECTORIES:%=test-%)

all: test

test: $(TEST_DIRECTORIES)

$(TEST_DIRECTORIES):
	$(MAKE) -C $(@:test-%=%) test

clean: $(CLEAN_DIRECTORIES)

$(CLEAN_DIRECTORIES):
	$(MAKE) -C $(@:clean-%=%) clean

.PHONY: test $(TEST_DIRECTORIES) $(CLEAN_DIRECTORIES)
