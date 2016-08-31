#!/bin/bash

set -e

if [ ! -z "$TRAVIS_TAG" ];then
	echo "the tag is $TRAVIS_TAG, will deploy...."
else
	echo "will not deploy..."
	exit 0
fi
cd ./cmd/server/
gox -output "../../dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

cd ../../cmd/client
gox -output "../../dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

cd ../../cmd/control
gox -output "../../dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

cd ../../
ghr $TRAVIS_TAG -u lovedboy -t $GITHUB_TOKEN -r gortcp  --replace  --debug  dist/
