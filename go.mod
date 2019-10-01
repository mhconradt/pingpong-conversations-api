module github.com/mhconradt/pingpong-conversations-api

require (
	github.com/apache/pulsar/pulsar-client-go v0.0.0-20190924204837-ee42cf403349
	github.com/golang/protobuf v1.3.2
	github.com/infobloxopen/atlas-app-toolkit v0.19.0
	github.com/infobloxopen/protoc-gen-gorm v0.18.0
	github.com/jinzhu/gorm v1.9.10
	github.com/micro/go-micro v1.10.0
	google.golang.org/genproto v0.0.0-20190916214212-f660b8655731
)

require github.com/mhconradt/grpc-statuses v0.0.0

replace github.com/mhconradt/grpc-statuses => ../grpc-statuses

require github.com/mhconradt/proto/status v0.0.0

replace github.com/mhconradt/proto/status => ../proto/status

require github.com/mhconradt/proto/conversation v0.0.0

require github.com/mhconradt/proto/message v0.0.0

replace github.com/mhconradt/proto/message => ../proto/message

replace github.com/mhconradt/proto/conversation => ../proto/conversation

require google.golang.org/grpc v1.22.1

require github.com/mhconradt/proto/user v0.0.0

replace github.com/mhconradt/proto/user => ../proto/user

go 1.13
