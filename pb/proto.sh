# /bin/bash

cd task && protoc  --go_out=plugins=grpc:. *.proto

echo "success"