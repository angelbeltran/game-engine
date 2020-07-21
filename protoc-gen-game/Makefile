default:
	protoc -I=protos --go_out=$$GOPATH/src game_engine.proto
	go mod tidy
	go mod vendor
	go build -o plugin *.go;
