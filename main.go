package main

import (
	"github.com/jinzhu/gorm"
	proto "github.com/mhconradt/pingpong-conversations-api/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server/grpc"
	"google.golang.org/grpc/encoding"
	_ "google.golang.org/grpc/encoding/proto"
	"log"
)

var pool *gorm.DB

func main() {
	/*
		Options:
		1. Move message and conversation into the same protobuf and golang packages in the proto repo.
		Pros: Allow has-many relationship between conversation and messages
		Cons: Kind of a pain in the ass to update the import statements, but really not all that bad.
		2. Write a raw SQL statement to get all of a user's conversations
		Pros: Clear. Can query on the status field of user_conversations.
		Cons: More work to move the things around and such.
	*/
	c := encoding.GetCodec("proto")
	codecOpt := grpc.Codec("application/protobuf", c)
	srv := grpc.NewServer(codecOpt)
	p := LookupWithDefault("PORT", ":3500")
	portOpt := micro.Address(p)
	// initialize database connection and store it in context...
	service := micro.NewService(
		micro.Server(srv),
		micro.Name("conversations"),
		micro.Version("latest"),
		portOpt,
	)
	var err error
	pool, err = InitDB()

	if err != nil {
		log.Fatal(err)
	}

	service.Init()

	if err = proto.RegisterConversationsHandler(service.Server(), new(Conversations)); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
