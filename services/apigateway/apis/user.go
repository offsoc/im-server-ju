package apis

import (
	"fmt"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/apigateway/models"
	"im-server/services/apigateway/services"
	"im-server/services/commonservices"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

func Register(ctx *gin.Context) {
	var userInfo models.UserInfo
	if err := ctx.BindJSON(&userInfo); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, resp, err := services.SyncApiCall(ctx, "reg_user", "", userInfo.UserId, &pbobjs.UserInfo{
		UserId:       userInfo.UserId,
		Nickname:     userInfo.Nickname,
		UserPortrait: userInfo.UserPortrait,
		ExtFields:    commonservices.Map2KvItems(userInfo.ExtFields),
	}, func() proto.Message {
		return &pbobjs.UserRegResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	rpcResp, ok := resp.(*pbobjs.UserRegResp)
	if !ok {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_RESP_FAIL)
		return
	}
	tools.SuccessHttpResp(ctx, models.UserRegResp{
		UserId: rpcResp.UserId,
		Token:  rpcResp.Token,
	})
}

func UpdateUser(ctx *gin.Context) {
	var req models.UserInfo
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "upd_user_info", "", req.UserId, &pbobjs.UserInfo{
		UserId:       req.UserId,
		Nickname:     req.Nickname,
		UserPortrait: req.UserPortrait,
		ExtFields:    commonservices.Map2KvItems(req.ExtFields),
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryUserInfo(ctx *gin.Context) {
	userid := ctx.Query("user_id")
	if userid == "" {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	code, userInfoObj, err := services.SyncApiCall(ctx, "qry_user_info", userid, userid, &pbobjs.UserIdReq{
		UserId: userid,
	}, func() proto.Message {
		return &pbobjs.UserInfo{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	userInfo := userInfoObj.(*pbobjs.UserInfo)
	tools.SuccessHttpResp(ctx, &models.UserInfo{
		UserId:       userInfo.UserId,
		Nickname:     userInfo.Nickname,
		UserPortrait: userInfo.UserPortrait,
		ExtFields:    commonservices.Kvitems2Map(userInfo.ExtFields),
		UpdatedTime:  userInfo.UpdatedTime,
	})
}

func KickUsers(ctx *gin.Context) {
	var req models.KickUserReq
	if err := ctx.BindJSON(&req); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	code, _, err := services.SyncApiCall(ctx, "kick_user", "", req.UserId, &pbobjs.KickUserReq{
		UserId:    req.UserId,
		Platforms: req.Platforms,
		DeviceIds: req.DeviceIds,
		Ext:       req.Ext,
	}, nil)
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}
	tools.SuccessHttpResp(ctx, nil)
}

func QryUserOnlineStatus(ctx *gin.Context) {
	var userOnlineReq models.UserOnlineStatusReq
	if err := ctx.BindJSON(&userOnlineReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}

	if len(userOnlineReq.UserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	ret := models.UserOnlineStatusResp{
		Items: []*models.UserOnlineStatusItem{},
	}
	tmpMap := sync.Map{}
	groups := bases.GroupTargets("qry_online_status", userOnlineReq.UserIds)
	wg := sync.WaitGroup{}
	for _, ids := range groups {
		wg.Add(1)
		userIds := ids
		go func() {
			defer wg.Done()
			_, resp, err := services.SyncApiCall(ctx, "qry_online_status", "", userIds[0], &pbobjs.UserOnlineStatusReq{
				UserIds: userIds,
			}, func() proto.Message {
				return &pbobjs.UserOnlineStatusResp{}
			})
			if err == nil {
				onlineResp, ok := resp.(*pbobjs.UserOnlineStatusResp)
				if ok && len(onlineResp.Items) > 0 {
					for _, item := range onlineResp.Items {
						tmpMap.Store(item.UserId, item)
					}
				}
			}
		}()
	}
	wg.Wait()
	tmpMap.Range(func(key, value any) bool {
		item := value.(*pbobjs.UserOnlineItem)
		ret.Items = append(ret.Items, &models.UserOnlineStatusItem{
			UserId:   item.UserId,
			IsOnline: item.IsOnline,
		})
		return true
	})
	tools.SuccessHttpResp(ctx, ret)
}

func UserBan(ctx *gin.Context) {
	var banReq models.BanUsersReq
	if err := ctx.BindJSON(&banReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(banReq.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	groups := map[string][]*pbobjs.BanUser{}
	for _, user := range banReq.Items {
		node := bases.GetCluster().GetTargetNode("ban_users", user.UserId)
		if node != nil && node.Name != "" {
			var arr []*pbobjs.BanUser
			var ok bool
			endTime := user.EndTime
			if user.BanType == int(pbobjs.BanType_Temporary) && endTime <= 0 {
				if user.EndTimeOffset > 0 {
					endTime = time.Now().UnixMilli() + user.EndTimeOffset
				}
			}
			pbBanUser := &pbobjs.BanUser{
				UserId:     user.UserId,
				BanType:    pbobjs.BanType(user.BanType),
				EndTime:    endTime,
				ScopeKey:   user.ScopeKey,
				ScopeValue: user.ScopeValue,
				Ext:        user.Ext,
			}
			if arr, ok = groups[node.Name]; ok {
				arr = append(arr, pbBanUser)
			} else {
				arr = []*pbobjs.BanUser{pbBanUser}
			}
			groups[node.Name] = arr
		}
	}
	wg := sync.WaitGroup{}
	for _, banUsers := range groups {
		wg.Add(1)
		users := banUsers
		go func() {
			defer wg.Done()
			services.AsyncApiCall(ctx, "ban_users", "", users[0].UserId, &pbobjs.BanUsersReq{
				BanUsers: users,
				IsAdd:    true,
			})
		}()
	}
	wg.Wait()
	tools.SuccessHttpResp(ctx, nil)
}

func UserUnBan(ctx *gin.Context) {
	var banReq models.BanUsersReq
	if err := ctx.BindJSON(&banReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(banReq.Items) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	groups := map[string][]*pbobjs.BanUser{}
	for _, user := range banReq.Items {
		node := bases.GetCluster().GetTargetNode("ban_users", user.UserId)
		if node != nil && node.Name != "" {
			var arr []*pbobjs.BanUser
			var ok bool
			pbBanUser := &pbobjs.BanUser{
				UserId:     user.UserId,
				BanType:    pbobjs.BanType(user.BanType),
				EndTime:    user.EndTime,
				ScopeKey:   user.ScopeKey,
				ScopeValue: user.ScopeValue,
			}
			if arr, ok = groups[node.Name]; ok {
				arr = append(arr, pbBanUser)
			} else {
				arr = []*pbobjs.BanUser{pbBanUser}
			}
			groups[node.Name] = arr
		}
	}
	wg := sync.WaitGroup{}
	for _, banUsers := range groups {
		wg.Add(1)
		users := banUsers
		go func() {
			defer wg.Done()
			services.AsyncApiCall(ctx, "ban_users", "", users[0].UserId, &pbobjs.BanUsersReq{
				BanUsers: users,
				IsAdd:    false,
			})
		}()
	}
	wg.Wait()
	tools.SuccessHttpResp(ctx, nil)
}

func QryBanUsers(ctx *gin.Context) {
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	offsetStr := ctx.Query("offset")

	code, resp, err := services.SyncApiCall(ctx, "qry_ban_users", "", fmt.Sprintf("%s%d", services.GetCtxString(ctx, services.CtxKey_AppKey), tools.RandInt(100000)), &pbobjs.QryBanUsersReq{
		Limit:  limit,
		Offset: offsetStr,
	}, func() proto.Message {
		return &pbobjs.QryBanUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	ret := &models.QryBanUsersResp{
		Items: []*models.BanUser{},
	}
	banUsers := resp.(*pbobjs.QryBanUsersResp)
	ret.Offset = banUsers.Offset
	for _, u := range banUsers.Items {
		ret.Items = append(ret.Items, &models.BanUser{
			UserId:      u.UserId,
			BanType:     int(u.BanType),
			CreatedTime: u.CreatedTime,
			EndTime:     u.EndTime,
			ScopeKey:    u.ScopeKey,
			ScopeValue:  u.ScopeValue,
			Ext:         u.Ext,
		})
	}
	tools.SuccessHttpResp(ctx, ret)
}

func BlockUser(ctx *gin.Context) {
	var blockReq models.BlockUsersReq
	if err := ctx.BindJSON(&blockReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(blockReq.BlockUserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	services.SyncApiCall(ctx, "block_users", blockReq.UserId, blockReq.UserId, &pbobjs.BlockUsersReq{
		UserIds: blockReq.BlockUserIds,
		IsAdd:   true,
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func UnBlockUser(ctx *gin.Context) {
	var blockReq models.BlockUsersReq
	if err := ctx.BindJSON(&blockReq); err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_REQ_BODY_ILLEGAL)
		return
	}
	if len(blockReq.BlockUserIds) <= 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_PARAM_REQUIRED)
		return
	}
	services.SyncApiCall(ctx, "block_users", blockReq.UserId, blockReq.UserId, &pbobjs.BlockUsersReq{
		UserIds: blockReq.BlockUserIds,
		IsAdd:   false,
	}, nil)
	tools.SuccessHttpResp(ctx, nil)
}

func QryBlockUsers(ctx *gin.Context) {
	userId := ctx.Query("user_id")
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	offsetStr := ctx.Query("offset")
	code, resp, err := services.SyncApiCall(ctx, "qry_block_users", userId, userId, &pbobjs.QryBlockUsersReq{
		Limit:  limit,
		Offset: offsetStr,
	}, func() proto.Message {
		return &pbobjs.QryBlockUsersResp{}
	})
	if err != nil {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode_API_INTERNAL_TIMEOUT)
		return
	}
	if code > 0 {
		tools.ErrorHttpResp(ctx, errs.IMErrorCode(code))
		return
	}

	ret := &models.QryBlockUsersResp{
		UserId: userId,
		Items:  []*models.BlockUser{},
	}
	blockUsers := resp.(*pbobjs.QryBlockUsersResp)
	ret.Offset = blockUsers.Offset
	for _, u := range blockUsers.Items {
		ret.Items = append(ret.Items, &models.BlockUser{
			BlockUserId: u.BlockUserId,
			CreatedTime: u.CreatedTime,
		})
	}
	tools.SuccessHttpResp(ctx, ret)
}
