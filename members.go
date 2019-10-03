package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/proto"
	status "github.com/mhconradt/grpc-statuses"
	proto2 "github.com/mhconradt/pingpong-conversations-api/proto"
	"github.com/mhconradt/pingpong-conversations-api/pubsub"
	conversation "github.com/mhconradt/proto/conversation"
	"github.com/mhconradt/proto/user"
)

// When you add someone, you need to tell them and the group

/*

	// Subscribes the user to a stream of the conversation's users
	// Lists existing users at the beginning of the session
	// Notifies the user of any new or removed users throughout the session
	ListMembers(context.Context, *ListMembersRequest, Conversations_ListMembersStream) error
	// Removes a user from the conversation
	RemoveMember(context.Context, *RemoveMemberRequest, *RemoveMemberResponse) error
*/

// todo: retrieve user with client instead of the database
// todo: send updates to conversation on the list conversations channel
// todo: send updates to profiles on the list members channel
// note: both of the above are practical assuming they will almost never happen. otherwise they are not.
// todo: system messages about updates to the conversation
// todo: system messages about removing or adding members
// todo: client
// todo: notifications

func (c *Conversations) AddMember(ctx context.Context, req *proto2.AddMemberRequest, res *proto2.AddMemberResponse) error {
	// Create the relation
	// This will not happen super frequently
	// Get conversation
	if err := saveUserStatus(ctx, req.UserId, req.ConversationId, 1); err != nil {
		res.Status = status.InternalServerError("failed to add member to conversation")
	}
	nuErr := notifyUser(ctx, req.UserId, req.ConversationId, proto2.Action_LIST)
	ngErr := notifyGroup(ctx, req.UserId, req.ConversationId, proto2.Action_LIST)
	if nuErr != nil {
		res.Status = status.InternalServerError("failed to invite new member")
		return nuErr
	}
	if ngErr != nil {
		res.Status = status.InternalServerError("failed to notify group")
		return ngErr
	}
	res.Status = status.SuccessValue
	return nil
}

func saveUserStatus(ctx context.Context, uId, cId int32, status int) error {
	q := fmt.Sprintf("insert into user_conversations (conversation_id, user_id, status) VALUES (%v, %v, %v) " +
		"ON CONFLICT ON CONSTRAINT user_conversations_pkey DO UPDATE SET status = %v;", cId, uId, status, status)
	stmt, err := pool.DB().PrepareContext(ctx, q)
	if err != nil {
		return err
	}
	defer func () { _ = stmt.Close() }()
	_, err = stmt.ExecContext(ctx)
	return err
}

func notifyUser(ctx context.Context, userId, conversationId int32, action proto2.Action) error {
	c, err := conversation.DefaultReadConversation(ctx, &conversation.Conversation{Id: conversationId}, pool)
	if err != nil {
		return err
	}
	msg, err := proto.Marshal(&proto2.ListConversationsResponse{
		Conversation: c,
		Action: action,
		Status: status.SuccessValue,
	})
	if err != nil {
		return err
	}
	topic := fmt.Sprintf("users_%v_conversations", userId)
	return pubsub.Publish(topic,  msg)
}

func notifyGroup(ctx context.Context, userId, conversationId int32, action proto2.Action) error {
	// This is hmm. Should change it to a client call.
	u, err := user.DefaultReadUser(ctx, &user.User{Id: userId}, pool)
	if err != nil {
		return err
	}
	topic := fmt.Sprintf("conversations_%v_members", conversationId)
	lmr := &proto2.ListMembersResponse{
		User: u,
		Action: action,
		Status: status.SuccessValue,
	}
	msg, err := proto.Marshal(lmr)
	if err != nil {
		return err
	}
	return pubsub.Publish(topic, msg)
}

func (c *Conversations) ListMembers(ctx context.Context, req *proto2.ListMembersRequest, s proto2.Conversations_ListMembersStream) error {
	members, err := fetchExistingMembers(ctx, req.ConversationId)
	if err != nil {
		res := &proto2.ListMembersResponse{
			Status: status.InternalServerError("failed to fetch members!"),
		}
		return s.Send(res)
	}
	if err := publishMembers(members, s); err != nil {
		res := &proto2.ListMembersResponse{
			Status: status.InternalServerError("failed to publish members"),
		}
		return s.Send(res)
	}
	topic := fmt.Sprintf("conversations_%v_members", req.ConversationId)
	sub := pubsub.NewSubscriber(topic)
	defer sub.Close()
	for msg := range sub.Messages {
		// send the list member response
		// send a message on add and remove
		lmr := new(proto2.ListMembersResponse)
		if err := proto.Unmarshal(msg, lmr); err != nil {
			res := &proto2.ListMembersResponse{
				Status: status.InternalServerError("failed to decode message"),
			}
			return s.Send(res)
		}
		_ = s.Send(lmr)
	}
	return nil
}

func fetchExistingMembers(ctx context.Context, conversationId int32) (members []*user.User, err error) {
	qs := fmt.Sprintf("with conversation_users as (" +
		"select * from conversations c inner join user_conversations uc on c.id = uc.conversation_id where c.id = %v) " +
		"select u.* from conversation_users inner join users u on u.id = conversation_users.user_id where conversation_users.status = 1;", conversationId)
	rows, err := pool.DB().QueryContext(ctx, qs)
	if err != nil {
		return members, err
	}
	return decodeMemberRows(ctx, rows)
}

func decodeMemberRows(ctx context.Context, rows *sql.Rows) ([]*user.User, error) {
	members := make([]*user.User, 0)
	for rows.Next() {
		u := new(user.UserORM)
		if err := pool.ScanRows(rows, u); err != nil {
			return members, err
		}
		pb, _ := u.ToPB(ctx)
		members = append(members, &pb)
	}
	return members, nil
}

func publishMembers(members []*user.User, s proto2.Conversations_ListMembersStream) error {
	for _, member := range members {
		res := &proto2.ListMembersResponse{
			User: member,
			Status: status.SuccessValue,
			Action: proto2.Action_LIST,
		}
		if err := s.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conversations) RemoveMember(ctx context.Context, req *proto2.RemoveMemberRequest, res *proto2.RemoveMemberResponse) error {
	if err := saveUserStatus(ctx, req.UserId, req.ConversationId, 0); err != nil {
		res.Status = status.InternalServerError("failed to update status")
		return err
	}
	// Notifying the user is an independent process
	// Save the status, then both notify the user and the group regardless of each other's success
	nuErr := notifyUser(ctx, req.UserId, req.ConversationId, proto2.Action_REMOVE)
	ngErr := notifyGroup(ctx, req.UserId, req.ConversationId, proto2.Action_REMOVE)
	if nuErr != nil {
		res.Status = status.InternalServerError("failed to notify the user")
		return nuErr
	}
	if ngErr != nil {
		res.Status = status.InternalServerError("failed to notify group")
		return ngErr
	}
	res.Status = status.SuccessValue
	// notify the user
	// notify the group
	return nil
}
