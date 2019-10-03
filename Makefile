protos:
	protoc --go_out=./proto/user/ --gorm_out=./proto/user/ -I=${GOPATH}/src -I=./proto/user/ ./proto/user/user.proto
	protoc --go_out=./proto/message/ --gorm_out=./proto/message/ -I=${GOPATH}/src -I=./proto/message/ ./proto/message/message.proto
	protoc --go_out=./proto/conversation/ --gorm_out=./proto/conversation/ -I=${GOPATH}/src -I=./proto/conversation/ ./proto/conversation/conversation.proto
	protoc --go_out=./proto/status/ --gorm_out=./proto/status/ -I=${GOPATH}/src -I=./proto/status/ ./proto/status/status.proto
