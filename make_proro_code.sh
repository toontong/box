#bin/sh

cmd='protoc --go_out=plugins=grpc:. *.proto'

proto_dirs="$PWD/proto"

for d in `ls $proto_dirs`
do
		cd "$proto_dirs/$d"
		$cmd
		cd -
done
