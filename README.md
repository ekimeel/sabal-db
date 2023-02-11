## sabal-db
sabal-db is the data store for the sabal project.

### protobuf 
protoc -I=./protobuf --go_out=. --go_opt=paths=source_relative \
--go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
protobuf/*.proto
