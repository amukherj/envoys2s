#!/usr/bin/env sh

if [ $# -lt 3 ]; then
	echo "usage: $0 <binary_name> <port> <envoy_config>"
	exit 1
fi

/usr/local/bin/$1 --port $2 &
/usr/local/bin/envoy -c $3
