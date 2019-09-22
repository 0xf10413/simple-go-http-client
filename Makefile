GOPATH=${PWD}
GOBIN=${GOPATH}/bin

launch-client:
	GOPATH=${GOPATH} go run client.go common.go

launch-server:
	env FLASK_DEBUG=1 FLASK_APP=server.py flask run -p 8081

deps:
	GOPATH=${GOPATH} GOBIN=${GOBIN} go get .
