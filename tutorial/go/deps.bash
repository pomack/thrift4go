#!/usr/bin/env bash
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

STARTDIR=`pwd`
GOLIBDIR=${GOROOT}/pkg/${GOOS}_${GOARCH}

OUT="Make.deps"
TMP="Make.deps.tmp"

if [ -f $OUT ] && ! [ -w $OUT ]; then
	echo "$0: $OUT is read-only; aborting." 1>&2
	exit 1
fi

# Get list of directories from Makefile
localdirs=$(gomake --no-print-directory echo-dirs)

dirs=""

function addGorootDep {
  FILENAME=$1
  FULLPATH=${GOROOT}/pkg/${GOOS}_${GOARCH}/${FILENAME}
  if [ -d ${FULLPATH} ]; then
    for filename in $(ls -1 ${FULLPATH}); do
      addGorootDep $1/$filename
    done
  else
    BASEFILENAME=$(basename $1)
    LIBNAME=$(basename -s .a $1)
    DIRNAME=$(dirname $1)
    if [ ${LIBNAME} != ${BASEFILENAME} ]; then
      if [ ${DIRNAME} = "." ]; then
        echo ${LIBNAME}.install: \${GOROOT}/pkg/\${GOOS}_\${GOARCH}/$1
        dirs="${dirs} ${LIBNAME}"
      else
        echo ${DIRNAME}/${LIBNAME}.install: \${GOROOT}/pkg/\${GOOS}_\${GOARCH}/$1
        dirs="${dirs} ${DIRNAME}/${LIBNAME}"
      fi
    fi
  fi
}

pushd . > /dev/null
(
for filename in $(ls -1 ${GOLIBDIR}); do
  addGorootDep $filename
done

dirpat=$(echo $dirs $localdirs C | awk '{
	for(i=1;i<=NF;i++){ 
		x=$i
		gsub("/", "\\/", x)
		printf("/^(%s)$/\n", x)
	}
}')

for dir in $localdirs; do 
  if [ -d $STARTDIR/$dir ]; then
    cd $STARTDIR/$dir || exit 1
  else
    continue
  fi

	sources=$(sed -n 's/^[ 	]*\([^ 	]*\.go\)[ 	]*\\*[ 	]*$/\1/p' Makefile)
	sources=$(echo $sources | sed 's/\$(GOOS)/'$GOOS'/g')
	sources=$(echo $sources | sed 's/\$(GOARCH)/'$GOARCH'/g')
	# /dev/null here means we get an empty dependency list if $sources is empty
	# instead of listing every file in the directory.
	sources=$(ls $sources /dev/null 2> /dev/null)  # remove .s, .c, etc.
	
	#sources=$(sed -n 's/\.go\\/.go/p' Makefile)
	#sources=$(ls $sources 2> /dev/null)  # remove .s, .c, etc.

	deps=$(
		sed -n '/^import.*"/p; /^import[ \t]*(/,/^)/p' $sources /dev/null |
		cut -d '"' -f2 |
		awk "$dirpat" |
		grep -v "^$dir\$" |
		sed 's/$/.install/' |
		sed 's;^C\.install;runtime/cgo.install;' |
		sort -u
	)

	echo $dir.install: $deps
done
) > $TMP

popd > /dev/null

mv $TMP $OUT
