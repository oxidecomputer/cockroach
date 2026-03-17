#!/usr/bin/env bash

set -euo pipefail

echo "The following types are considered always safe for reporting:"
echo
echo "File | Type"; echo "--|--"
grep -r -n '^func \(.*\) SafeValue\(\)' | \
    grep -v '^vendor/github.com/cockroachdb/redact' | \
    sed -E -e 's/^([^:]*):[0-9]+:func \(([^ ]* )?(.*)\) SafeValue.*$$/\1 | \`\3\`/g' | \
    LC_ALL=C sort
grep -r -n 'redact\.RegisterSafeType' | \
    grep -vE '^([^:]*):[0-9]+:[ 	]*//' | \
    grep -v '^vendor/github.com/cockroachdb/redact' | \
    sed -E -e 's/^([^:]*):[0-9]+:.*redact\.RegisterSafeType\((.*)\).*/\1 | \`\2\`/g' | \
    LC_ALL=C sort
