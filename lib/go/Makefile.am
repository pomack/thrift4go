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

all: test

export GOPATH = $(PWD)
SOURCE = $(GOPATH)/{pkg,src}
DESTINATION = $(DESTDIR)/$(GOROOT)

clean-local:
	go clean thrift

nuke-local:
	go clean -i -r -x thrift
	$(RM) -rf bin pkg

# N.B.(matttproud): This step may fail to run with sudo due to sudo clearing
#                   environment variables before changing permissions for
#                   security reasons.  This can be avoided by using the
#                   DESTDIR construct described above to provide emissions to
#                   an intermediate directory for subsequent manual merging.

install-exec-local: test
	go install -v -x thrift
	install -d "$(GOPATH)/pkg"
	install -d "$(GOPATH)/src"
	cp -av "$(GOPATH)/pkg/" "$(DESTINATION)/pkg/"
	cp -av "$(GOPATH)/src/" "$(DESTINATION)/src/"

test: check
	go test -v -x thrift

check-local:
	go build -v -x thrift
