protoc --go_out=. --go_opt=paths=source_relative --proto_path=. protos.proto
protoc --go_out=. --go_opt=paths=source_relative --plugin=$GOPATH/bin/protoc-gen-go-grain --go-grain_out=. --go-grain_opt=paths=source_relative protos.proto
