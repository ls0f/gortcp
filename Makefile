BIN := server client control


test:
	go test -cover  `go list ./...`

fmt:
	find . -name "*.go" -type f -exec echo {} \; | grep -v -E "github.com|gopkg.in"|\
	while IFS= read -r line; \
	do \
		echo "$$line";\
		goimports -w "$$line" "$$line";\
	done

build:
	echo ==================================; \
	for m in $(BIN); do \
		cd $(PWD)/cmd/$$m && go build --race -o $$m $$m.go; \
	done
	echo ==================================; \
