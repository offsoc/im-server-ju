package friends

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/gmicro/actorsystem"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/appbusiness/services"

	"google.golang.org/protobuf/proto"
)

type QryFriendsActor struct {
	bases.BaseActor
}

func (actor *QryFriendsActor) OnReceive(ctx context.Context, input proto.Message) {
	if req, ok := input.(*pbobjs.FriendListReq); ok {
		code, resp := services.QryFriends(ctx, req)
		ack := bases.CreateQueryAckWraper(ctx, code, resp)
		actor.Sender.Tell(ack, actorsystem.NoSender)
	}
}

func (actor *QryFriendsActor) CreateInputObj() proto.Message {
	return &pbobjs.FriendListReq{}
}
