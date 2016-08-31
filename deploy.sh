#!/bin/bash


if [ "$TRAVIS_BRANCH" = "master" ] && [ ! -z "$TRAVIS_TAG" ];then
    echo "This will deploy!"
  else
    echo "This will not deploy!"
    exit 0
fi

go get github.com/mitchellh/gox
go get github.com/tcnksm/ghr
  
cd ./cmd/server/
gox -output "../../dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

cd ../../cmd/client
gox -output "../../dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

cd ../../cmd/control
gox -output "../../dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

cd ../../
ghr $TRAVIS_TAG --username lovedboy --token $GITHUB_TOKEN  --replace  --debug  dist/
