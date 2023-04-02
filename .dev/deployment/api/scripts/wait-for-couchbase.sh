#!/bin/sh

set -e

host="$1"
shift
cmd="$@"


until curl -sSf -u "admin:password" http://couchbase:8091/pools/default > /dev/null; do
  sleep 1
done

sleep 10

echo "Couchbase is ready!"

exec ./app
