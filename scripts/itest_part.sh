#!/bin/bash

# Let's work with absolute paths only, we run in the itest directory itself.
WORKDIR=$(pwd)/lntest/itest

TRANCHE=$1
NUM_TRANCHES=$2

# Shift the passed parameters by two, giving us all remaining testing flags in
# the $@ special variable.
shift
shift

# Windows insists on having the .exe suffix for an executable, we need to add
# that here if necessary.
EXEC="$WORKDIR"/itest.test"$EXEC_SUFFIX"
broln_EXEC="$WORKDIR"/broln-itest"$EXEC_SUFFIX"
brond_EXEC="$WORKDIR"/brond-itest"$EXEC_SUFFIX"
echo $EXEC -test.v "$@" -logoutput -goroutinedump -logdir=.logs-tranche$TRANCHE -brolnexec=$broln_EXEC -brondexec=$brond_EXEC -splittranches=$NUM_TRANCHES -runtranche=$TRANCHE

# Exit code 255 causes the parallel jobs to abort, so if one part fails the
# other is aborted too.
cd "$WORKDIR" || exit 255
$EXEC -test.v "$@" -logoutput -goroutinedump -logdir=.logs-tranche$TRANCHE -brolnexec=$broln_EXEC -brondexec=$brond_EXEC -splittranches=$NUM_TRANCHES -runtranche=$TRANCHE || exit 255
