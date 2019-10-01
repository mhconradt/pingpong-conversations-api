package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	status "github.com/mhconradt/grpc-statuses"
	proto2 "github.com/mhconradt/pingpong-conversations-api/proto"
	"github.com/mhconradt/pingpong-conversations-api/pubsub"
	message "github.com/mhconradt/proto/message"
)

func (c *Conversations) ListMessages(ctx context.Context, req *proto2.ListMessagesRequest, s proto2.Conversations_ListMessagesStream) error {
	messages, err := fetchExistingMessages(ctx, req.ConversationId, req.StartMessageId)
	if err != nil {
		res := &proto2.ListMessagesResponse{
			Status: status.InternalServerError("failed to decode query results"),
		}
		return s.Send(res)
	}
	if err := publishMessages(messages, s); err != nil {
		res := &proto2.ListMessagesResponse{
			Status: status.InternalServerError("failed to publish message!"),
		}
		return s.Send(res)
	}
	topicName := fmt.Sprintf("conversations_%v_messages", req.ConversationId)
	sub, err := pubsub.NewSubscriber(topicName, fmt.Sprintf("%v", req.SubscriberId))
	if err != nil {
		res := &proto2.ListMessagesResponse{
			Status: status.InternalServerError("failed to subscribe to new messages for conversation: ", req.ConversationId),
		}
		return s.Send(res)
	}
	defer sub.Close()
	for msg := range sub.Messages {
		m := new(message.Message)
		if err := proto.Unmarshal(msg, m); err != nil {
			res := &proto2.ListMessagesResponse{
				Status: status.InternalServerError("failed to decode message"),
			}
			return s.Send(res)
		}
		res := &proto2.ListMessagesResponse{
			Message: m,
			Status:  status.SuccessValue,
		}
		_ = s.Send(res)
	}
	return nil
}

func fetchExistingMessages(ctx context.Context, conversationId, startId int32) (messages []*message.Message, err error) {
	// do the fetch
	// decode the rows
	// return the messages
	rows, err := pool.DB().QueryContext(ctx,
		fmt.Sprintf("with conversation_messages as ("+
			"select m.* from conversations c inner join messages m on c.id = m.conversation_id where c.id = %v"+
			") select * from conversation_messages where id > %v;", conversationId, startId))
	if err != nil {
		return messages, err
	}
	return decodeMessageRows(ctx, rows)
}

func decodeMessageRows(ctx context.Context, rows *sql.Rows) ([]*message.Message, error) {
	messages := make([]*message.Message, 0)
	for rows.Next() {
		m := new(message.MessageORM)
		if err := pool.ScanRows(rows, m); err != nil {
			return messages, err
		}
		pb, _ := m.ToPB(ctx)
		messages = append(messages, &pb)
	}
	return messages, nil
}

func publishMessages(messages []*message.Message, s proto2.Conversations_ListMessagesStream) error {
	for _, msg := range messages {
		res := &proto2.ListMessagesResponse{
			Message: msg,
			Status:  status.SuccessValue,
		}
		// send latency may be non-negligible
		if err := s.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conversations) SendMessage(ctx context.Context, req *proto2.SendMessageRequest, res *proto2.SendMessageResponse) error {
	// write to db
	m, err := message.DefaultCreateMessage(ctx, req.Message, pool)
	if err != nil {
		res.Status = status.InternalServerError("failed to create message")
		return err
	}
	res.Message = m
	msg, err := proto.Marshal(m)
	if err != nil {
		res.Status = status.InternalServerError("failed to encode message before publishing")
		return err
	}
	topicName := fmt.Sprintf("conversations_%v_messages", m.ConversationId)
	if err := pubsub.SendOnce(topicName, msg); err != nil {
		res.Status = status.InternalServerError("failed to publish message")
		return err
	}
	res.Status = status.SuccessValue
	// send to Pulsar
	return nil
}
