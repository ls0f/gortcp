#!/bin/bash

find . -name "*.go" -type f -exec echo {} \; | 

while IFS= read -r line
do
	echo $line
	goimports -w $line $line
done
