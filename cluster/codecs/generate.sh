#!/usr/bin/env bash

# Copyright (C) 2020-2022 Talos, Inc.

# This isn't that useful without modification or an identical layout to the
# defaults below so you might consider it more documentation than working
# script.

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
AERON_DIR=$DIR/../../../aeron
SBE_DIR=$DIR/../../../simple-binary-encoding
AGRONA_DIR=$DIR/../../../agrona
AGRONA_VERSION=1.12.0
# AGRONA_JAR=`find ~/.gradle/ | grep agrona-$AGRONA_VERSION.jar | head -1`
AGRONA_JAR=$AGRONA_DIR/agrona/build/libs/agrona-1.12.0.jar

if [ ! -d "$AERON_DIR" ]
then
    echo "Can't find aeron directory ($AERON_DIR)"
    exit 1
else
    echo "Using aeron from $AERON_DIR"
fi

if [ ! -d "$SBE_DIR" ]
then
    echo "Can't find SBE directory ($SBE_DIR)"
    exit 1
else
    echo "Using SBE from $SBE_DIR"
fi

if [ ! -f "$AGRONA_JAR" ]
then
    echo "Can't find agrona jar (looking for version $AGRONA_VERSION)"
    exit 1
else
    echo "Using Agrona from $AGRONA_JAR"
fi
 
java \
  -Dsbe.output.dir=$DIR/.. \
  -Dsbe.target.language=golang \
  -Dsbe.validation.xsd=$SBE_DIR/sbe-tool/src/main/resources/fpl/sbe.xsd \
  -Dfile.encoding=UTF-8 \
  -Dsbe.target.namespace=codecs \
  -cp $SBE_DIR/sbe-tool/build/classes/java/main:$SBE_DIR/sbe-tool/build/resources/main:$AGRONA_JAR \
  uk.co.real_logic.sbe.SbeTool \
  $AERON_DIR/aeron-cluster/src/main/resources/cluster/aeron-cluster-codecs.xml
