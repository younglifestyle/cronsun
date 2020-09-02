#!/bin/sh

# shellcheck disable=SC2112
function check_code() {
	EXCODE=$?
	if [ "$EXCODE" != "0" ]; then
		echo "build fail."
		exit $EXCODE
	fi
}

out="dist"
echo "build file to ./$out"

mkdir -p "$out/conf"

go build -o ./$out/cronnode ./bin/node/server.go
check_code
go build -o ./$out/cronsrv ./bin/web/server.go
check_code

sources=`find ./conf/files -name "*.json.sample"`
check_code
for source in $sources;do
	yes | echo $source|sed "s/.*\/\(.*\.json\).*/cp -f & .\/$out\/conf\/\1/"|bash
	check_code
done

cp control ./$out/

echo "build success."
