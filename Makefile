

.PHONY: cleanup
cleanup: ## Clear the working area and the project
	rm -rf ${HOME}/.config/gp
	rm ./gp


.PHONY: test
test: 
	go test -v -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html


.PHONY: build
build: 
	go build -o ./dist/g -ldflags="-X github.com/ipedrazas/gp/cmd.Version=`git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*'` -X github.com/ipedrazas/gp/cmd.CommitHash=`git rev-parse --short HEAD` -X github.com/ipedrazas/gp/cmd.BuildTimestamp=`date '+%Y-%m-%dT%H:%M:%S'`"
	
.PHONY: install
install:
	rm ~/go/bin/g
	cp ./dist/g ~/go/bin/g

