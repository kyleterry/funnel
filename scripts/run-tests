#!/usr/bin/env bash

list=$(go list ./... | grep -v '^github.com/kyleterry/funnel/vendor/')

echo "Extreme vetting..."
echo ${list} | xargs -n1 go vet
echo "Extreme testing..."
echo ${list} | xargs -n1 go test ${FUNNEL_TEST_FLAGS:--cover -timeout=360s}
echo "Extreme errchecking..."
echo ${list} | xargs -n1 errcheck
