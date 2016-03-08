cd ./nameserver/
protoc --go_out=plugins=grpc:. *.proto
cd ..

cd ./worker/
protoc --go_out=plugins=grpc:. *.proto
cd ..

