package pb

//go:generate protoc -I $GOPATH/src -I ..\..\..\proto --go_out=plugins=grpc,Mgithub.com/golang/protobuf/ptypes/timestamp/timestamp.proto=/github.com/golang/protobuf/ptypes/timestamp:. ..\..\..\proto\contentservice.proto
