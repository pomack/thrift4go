#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements. See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership. The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License. You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied. See the License for the
# specific language governing permissions and limitations
# under the License.
#

# N.B.(matttproud): The flag assumptions of this makefile could be improved to
#                   potentially be more cross-platform compatible.  For
#                   instance, I think the »-t« flag for »cp« may be invalid on
#                   BSD and Darwin systems.

all: install


export GOPATH = $(PWD)
SOURCE = $(GOPATH)/{pkg,/src}
DESTINATION = $(DESTDIR)/$(GOROOT)


clean-local:
	go clean thrift

nuke-local:
	go clean -i -r -x thrift
	$(RM) -rf bin pkg

install-exec-local: test
	go install -a -v -x thrift
	mkdir -vp "$(DESTINATION)"
	cp -av $(SOURCE) -t "$(DESTINATION)"

test: check
	go test -a -v -x thrift

check-local:
	go build -a -v -x thrift
