#!/bin/bash

if [ $# -lt 4 ]; then
    echo Usage: json2go schema.json schema.go package varname
    exit;
fi

SCHEMAFILE=$1
GOFILE=$2
PACK=$3
VAR=$4

echo "package ${PACK}" > $GOFILE
echo "" >> $GOFILE
echo "// ${VAR} was generated from ${SCHEMAFILE} at" `date`>> $GOFILE
echo "var ${VAR} = []byte(\`" >> $GOFILE
cat ${SCHEMAFILE} >> $GOFILE
echo "\`)" >> $GOFILE
