GOPATH=${PWD}
GOBIN=${GOPATH}/bin

launch-client:
	go run client.go common.go

launch-server:
	go run server.go common.go

deps:
	go get .
