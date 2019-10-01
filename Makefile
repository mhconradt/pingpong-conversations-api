protos:
	protoc --go_out=./proto/ --gorm_out=./proto/ --micro_out=./proto/ -I=${GOPATH}/src -I=./proto/ ./proto/conversations.proto
run:
	go build
	./pingpong-conversations-api
