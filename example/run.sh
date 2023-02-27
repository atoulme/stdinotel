#!/bin/bash
set -ex pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
GOOS=`go env GOOS`
GOARCH=`go env GOARCH`

export STDINOTEL_PROTOCOL=splunk_hec
export STDINOTEL_TOKEN=00000000-0000-0000-0000-0000000000000
export STDINOTEL_ENDPOINT=https://localhost:18088
export STDINOTEL_SPLUNK_INDEX=main
export STDINOTEL_TLS_INSECURE_SKIP_VERIFY=true


echo "foo" | $SCRIPT_DIR/../bin/stdinotel_${GOOS}_${GOARCH}
cat $SCRIPT_DIR/lorem.txt | $SCRIPT_DIR/../bin/stdinotel_${GOOS}_${GOARCH}