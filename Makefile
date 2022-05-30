

.PHONY: clear
clear: ## Clear the working area and the project
	rm -rf ${HOME}/.config/gp
	rm ./gp


.PHONY: test
test: 
	go test -v -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html


.PHONY: build
build: 
	go build -ldflags "-X github.com/ipedrazas/gp/cmd.Version=`git describe --match 'v[0-9]*' --dirty='.m' --always --tags` -X github.com/ipedrazas/gp/cmd.Sha1=`git rev-parse HEAD`" -o g

