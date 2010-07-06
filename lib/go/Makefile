
all: install

DIRS=\
     thrift/\
     thrift/transport\
     thrift/protocol\
     thrift/base\
     thrift/server\

TEST=\
    thrift/transport\
    thrift/protocol\
    thrift/base\
    thrift/server\

TEST=\
	$(filter-out $(NOTEST),$(DIRS))


clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
nuke.dirs: $(addsuffix .nuke, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))

%.clean:
	+cd $* && gomake clean

%.install:
	+cd $* && gomake install

%.nuke:
	+cd $* && gomake nuke

%.test:
	+cd $* && gomake test

clean: clean.dirs

install: install.dirs

test:	test.dirs

#nuke: nuke.dirs
#	rm -rf "$(GOROOT)"/pkg/*

deps:
	./deps.bash

-include Make.deps
