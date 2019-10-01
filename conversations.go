package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/mhconradt/grpc-statuses"
	proto2 "github.com/mhconradt/pingpong-conversations-api/proto"
	"github.com/mhconradt/pingpong-conversations-api/pubsub"
	conversation "github.com/mhconradt/proto/conversation"
)

type Conversations struct{}

func (c *Conversations) CreateConversation(ctx context.Context, req *proto2.CreateConversationRequest, res *proto2.CreateConversationResponse) error {
	// don't allow adding on create yet
	if con, err := conversation.DefaultCreateConversation(ctx, req.Conversation, pool); err != nil {
		res.Status = status.InternalServerError("failed to create con")
	} else {
		res.Conversation = con
		res.Status = status.Success()
	}
	return nil
}

func (c *Conversations) UpdateConversation(ctx context.Context, req *proto2.UpdateConversationRequest, res *proto2.UpdateConversationResponse) error {
	if con, err := conversation.DefaultPatchConversation(ctx, req.Conversation, req.UpdateMask, pool); err != nil {
		res.Status = status.BadRequest("failed to update conversation")
	} else {
		res.Conversation = con
		res.Status = status.Success()
	}
	return nil
}

func (c *Conversations) ListConversations(ctx context.Context, req *proto2.ListConversationsRequest, s proto2.Conversations_ListConversationsStream) error {
	// select the user
	// select the users conversations where uc.status = 1
	conversations, err := fetchExistingConversations(ctx, req.SubscriberId)
	if err != nil {
		res := &proto2.ListConversationsResponse{
			Status: status.InternalServerError("failed to fetch conversations"),
		}
		return s.Send(res)
	}
	if err := publishConversations(conversations, s); err != nil {
		res := &proto2.ListConversationsResponse{
			Status: status.InternalServerError("failed to publish conversations"),
		}
		return s.Send(res)
	}
	topic := fmt.Sprintf("users_%v_conversations", req.SubscriberId)
	sub, err := pubsub.NewSubscriber(topic, "self")
	if err != nil {
		res := proto2.ListConversationsResponse{
			Status: status.InternalServerError("error subscribing to new conversations"),
		}
		return s.Send(&res)
	}
	defer sub.Close()
	for msg := range sub.Messages {
		lcr := new(proto2.ListConversationsResponse)
		if err := proto.Unmarshal(msg, lcr); err != nil {
			res := proto2.ListConversationsResponse{
				Status: status.InternalServerError("failed to unmarshal conversation message"),
			}
			return s.Send(&res)
		}
		if err := s.Send(lcr); err != nil {
			res := proto2.ListConversationsResponse{
				Status: status.InternalServerError("failed to send notification to client"),
			}
			return s.Send(&res)
		}
		// This should be a conversation record, as well as the relevant ADD, REMOVE, etc.
	}
	return nil
}

func fetchExistingConversations(ctx context.Context, userId int32) (conversations []*conversation.Conversation, err error) {
	rows, err := pool.DB().QueryContext(ctx,
		fmt.Sprintf("with u as (select * from users where id = %v) select c.* from conversations c inner join user_conversations uc on c.id = uc.conversation_id where uc.status = 1;", userId))
	if err != nil {
		return conversations, err
	}
	return decodeConversationRows(ctx, rows)
}

func decodeConversationRows(ctx context.Context, rows *sql.Rows) ([]*conversation.Conversation, error) {
	conversations := make([]*conversation.Conversation, 0)
	for rows.Next() {
		co := new(conversation.ConversationORM)
		if err := pool.ScanRows(rows, co); err != nil {
			return conversations, err
		}
		pb, _ := co.ToPB(ctx)
		conversations = append(conversations, &pb)
	}
	return conversations, nil
}

func publishConversations(conversations []*conversation.Conversation, s proto2.Conversations_ListConversationsStream) error {
	for _, c := range conversations {
		res := &proto2.ListConversationsResponse{
			Conversation: c,
			Status: status.SuccessValue,
			Action: proto2.Action_LIST,
		}
		if err := s.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conversations) DeleteConversation(ctx context.Context, req *proto2.DeleteConversationRequest, res *proto2.DeleteConversationResponse) error {
	co := &conversation.Conversation{
		Id: req.Conversation.Id,
	}
	if err := conversation.DefaultDeleteConversation(ctx, co, pool); err != nil {
		res.Status = status.BadRequest("failed to delete co.")
	} else {
		res.Status = status.Success()
	}
	return nil
}
