#!/bin/bash

cat >&2 <<EOF
This script will patch in the local changes such that one can build Apache
Thrift 0.8.0 with the current revision of thrift4go.  A number of sanity
checks will be run before-hand.

One may wish to run this script as "bash -x ${0}" to see exactly what steps
it performs, as it will modify the local source copy of Thrift.  The script
will abort if it encounters any errors.
EOF

while true ; do
  read -p "Do you wish to proceed: [Y]es or [N]o?" response
  case "${response}" in
    [Yy]* ) break ;;
    [Nn]* ) echo "Aborting..." ; exit 1 ;;
  esac
done

if [ ! -d "${THRIFT}" ]; then
cat >&2 <<EOF
A required environmental variable was not set: "THRIFT".  This should refer to
the path where a pristine copy of Thrift 0.8.0 lives.

export THRIFT='/tmp/thrift-0.8.0'
EOF

exit 1
fi

if [ ! -d "${THRIFT4GO}" ]; then
cat >&2 <<EOF
A required environmental variable was not set: "THRIFT4GO".  This should refer
to the path where a pristine copy of thrift4g0 lives.

export THRIFT4GO='/tmp/thrift4go
EOF

exit 1
fi

set -e

rm -rf "${THRIFT}/lib/go"

cp -a "${THRIFT4GO}/lib/go" "${THRIFT}/lib/go"

cp "${THRIFT4GO}/lib/Makefile.am" "${THRIFT}/lib/Makefile.am"

cp "${THRIFT4GO}/configure.ac" "${THRIFT}/configure.ac"

cp "${THRIFT4GO}/compiler/cpp/src/generate/t_go_generator.cc" \
  "${THRIFT}/compiler/cpp/src/generate/t_go_generator.cc"

cd "${THRIFT}"

./cleanup.sh

./bootstrap.sh

echo "The Thrift package in ${THRIFT} is ready for configuration now." >&2
