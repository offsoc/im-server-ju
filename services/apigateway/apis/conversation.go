package apis

import (
	"fmt"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func AddConversation(ctx *gin.Context) {
	var req models.Conversation
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	services.SyncApiCall(ctx, "add_conver", req.UserId, req.UserId, &pbobjs.Conversation{
		UserId:      req.UserId,
		TargetId:    req.TargetId,
		ChannelType: pbobjs.ChannelType(req.ChannelType),
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func DelConversation(ctx *gin.Context) {
	var req models.Conversations
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(req.Items) <= 0 || req.UserId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	convers := []*pbobjs.Conversation{}
	for _, c := range req.Items {
		convers = append(convers, &pbobjs.Conversation{
			UserId:      c.UserId,
			TargetId:    c.TargetId,
			ChannelType: pbobjs.ChannelType(c.ChannelType),
		})
	}

	services.SyncApiCall(ctx, "del_convers", req.UserId, req.UserId, &pbobjs.ConversationsReq{
		Conversations: convers,
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func ClearConverUnread(ctx *gin.Context) {
	var req models.Conversations
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(req.Items) <= 0 || req.UserId == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	convers := []*pbobjs.Conversation{}
	for _, c := range req.Items {
		convers = append(convers, &pbobjs.Conversation{
			UserId:      req.UserId,
			TargetId:    c.TargetId,
			ChannelType: pbobjs.ChannelType(c.ChannelType),
		})
	}
	services.AsyncApiCall(ctx, "clear_unread", req.UserId, req.UserId, &pbobjs.ClearUnreadReq{
		Conversations: convers,
		NoCmdMsg:      true,
	})
	tools.SuccessHttpResp(ctx, nil)
}

// undisturb_convers
func UndisturbConvers(ctx *gin.Context) {
	var undisturbConversReq models.UndisturbConversReq
	if err := ctx.BindJSON(&undisturbConversReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(undisturbConversReq.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	items := []*pbobjs.UndisturbConverItem{}
	for _, reqItem := range undisturbConversReq.Items {
		items = append(items, &pbobjs.UndisturbConverItem{
			TargetId:      reqItem.TargetId,
			ChannelType:   pbobjs.ChannelType(reqItem.ChannelType),
			UndisturbType: reqItem.UndisturbType,
		})
	}
	services.SyncApiCall(ctx, "undisturb_convers", undisturbConversReq.UserId, undisturbConversReq.UserId, &pbobjs.UndisturbConversReq{
		Items: items,
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func QryGlobalConvers(ctx *gin.Context) {
	start := ctx.Query("start")
	var startTime int64 = 0
	if start != "" {
		intVal, err := tools.String2Int64(start)
		if err == nil {
			startTime = intVal
		}
	}

	limitStr := ctx.Query("count")
	var limit int64 = 100
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	rpcTargetId := fmt.Sprintf("random%d", tools.RandInt(1000))
	//targetId
	targetId := ctx.Query("target_id")
	//channelType
	channelTypeStr := ctx.Query("channel_type")
	channelTypeInt, err := tools.String2Int64(channelTypeStr)
	channelType := pbobjs.ChannelType_Unknown
	if err == nil {
		channelType = pbobjs.ChannelType(channelTypeInt)
	}
	//exclude user_ids
	excludeUserIds := ctx.QueryArray("exclude_user_id")
	if len(excludeUserIds) == 0 {
		excludeUserIds = ctx.QueryArray("exclude_user_ids")
	}
	code, resp, err := services.SyncApiCall(ctx, "qry_global_convers", "", rpcTargetId, &pbobjs.QryGlobalConversReq{
		Start:          startTime,
		Order:          0,
		Count:          int32(limit),
		TargetId:       targetId,
		ChannelType:    channelType,
		ExcludeUserIds: excludeUserIds,
	}, func() proto.Message {
		return &pbobjs.QryGlobalConversResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code != int32(errs.IMErrorCode_SUCCESS) {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	conversResp := resp.(*pbobjs.QryGlobalConversResp)
	ret := models.Conversations{
		Items:      []*models.Conversation{},
		IsFinished: conversResp.IsFinished,
	}
	for _, conver := range conversResp.Convers {
		item := &models.Conversation{
			Id:          conver.Id,
			UserId:      conver.SenderId,
			TargetId:    conver.TargetId,
			ChannelType: int(conver.ChannelType),
			Time:        conver.UpdatedTime,
		}
		ret.Items = append(ret.Items, item)
	}

	tools.SuccessHttpResp(ctx, ret)
}
