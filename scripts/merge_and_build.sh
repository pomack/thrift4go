#!/bin/bash

cat >&2 <<EOF
This script will patch in the local changes such that one can build Apache
Thrift 0.8.0 with the current revision of thrift4go.  A number of sanity
checks will be run before-hand.

One may wish to run this script as "bash -x ${0}" to see exactly what steps
it performs, as it will modify the local source copy of Thrift.  The script
will abort if it encounters any errors.

It may also be run as "${0} -b" for batch mode, in which case no prompts will
be performed.
EOF

batch=false

arguments=$(getopt b ${*})
set -- ${arguments}

for i do
  case "${i}" in
    -b)
      batch=true
      shift;;
    --)
      shift
      break;;
  esac
done

if [[ "${batch}" == "false" ]]; then
while true ; do
  read -p "Do you wish to proceed: [Y]es or [N]o?" response
  case "${response}" in
    [Yy]* ) break ;;
    [Nn]* ) echo "Aborting..." ; exit 1 ;;
  esac
done
fi

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
to the path where a pristine copy of thrift4go lives.

export THRIFT4GO='/tmp/thrift4go
EOF

exit 1
fi

if [ ! -x "${THRIFT}/cleanup.sh" ]; then
cat >&2 <<EOF
WARNING: ${THRIFT}/cleanup.sh does not exist.  A copy from Thrift 0.8.0 is
being provided in its place.  This may occur if using a Thrift release tarball.
EOF
cp -f "${THRIFT4GO}/scripts/cleanup.sh" "${THRIFT}/cleanup.sh"
fi

if [ ! -x "${THRIFT}/bootstrap.sh" ]; then
cat >&2 <<EOF
WARNING: ${THRIFT}/bootstrap.sh does not exist.  A copy from Thrift 0.8.0 is
being provided in its place.  This may occur if using a Thrift release tarball.
EOF
cp -f "${THRIFT4GO}/scripts/bootstrap.sh" "${THRIFT}/bootstrap.sh"
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
