#!/usr/bin/env bash -e

cd $(dirname $0)

for proto in *.proto
do
  echo "generating Go for $proto"
  protoc \
    -I . \
    --go_out=../pb \
    --go-grpc_out=../pb \
    --grpc-gateway_out=../pb \
    --openapiv2_out=../../../../doc   \
    --openapiv2_opt=disable_default_errors=true \
    $proto
done

rm ../../../../doc/error.swagger.json
