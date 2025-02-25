package services

import (
	"context"
	"im-server/commons/bases"
	"im-server/commons/errs"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/commons/tools"
	"im-server/services/appbusiness/models"
	"im-server/services/appbusiness/storages"
	storeModels "im-server/services/appbusiness/storages/models"
	"im-server/services/commonservices"
	"im-server/services/commonservices/msgdefines"
	"time"

	"google.golang.org/protobuf/proto"
)

func QryGroupInfo(ctx context.Context, groupId string) (errs.IMErrorCode, *pbobjs.GrpInfo) {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, respObj, err := bases.SyncRpcCall(ctx, "qry_group_info", groupId, &pbobjs.GroupInfoReq{
		GroupId: groupId,
	}, func() proto.Message {
		return &pbobjs.GroupInfo{}
	})
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	grpInfo := respObj.(*pbobjs.GroupInfo)
	ret := &pbobjs.GrpInfo{
		GroupId:       grpInfo.GroupId,
		GroupName:     grpInfo.GroupName,
		GroupPortrait: grpInfo.GroupPortrait,
		Members:       []*pbobjs.GroupMemberInfo{},
		MemberCount:   grpInfo.MemberCount,
		GroupManagement: &pbobjs.GroupManagement{
			GroupId:       grpInfo.GroupId,
			GroupMute:     grpInfo.IsMute,
			MaxAdminCount: 10,
		},
	}
	isMember := false
	//check is member
	code, respObj, err = bases.SyncRpcCall(ctx, "g_check_members", groupId, &pbobjs.CheckGroupMembersReq{
		GroupId:   groupId,
		MemberIds: []string{requestId},
	}, func() proto.Message {
		return &pbobjs.CheckGroupMembersResp{}
	})
	if err == nil && code == errs.IMErrorCode_SUCCESS && respObj != nil {
		checkResp, ok := respObj.(*pbobjs.CheckGroupMembersResp)
		if ok {
			if _, exist := checkResp.MemberIdMap[requestId]; exist {
				isMember = true
			}
		}
	}
	if !isMember {
		ret.MyRole = pbobjs.GrpMemberRole_GrpNotMember
		return errs.IMErrorCode_SUCCESS, ret
	}

	var creator string = ""
	administrators := map[string]bool{}
	for _, setting := range grpInfo.Settings {
		if setting.Key == string(commonservices.AttItemKey_GrpCreator) {
			creator = setting.Value
		} else if setting.Key == string(commonservices.AttItemKey_GrpAdministrators) {
			if len(setting.Value) > 0 {
				adminIds := []string{}
				err := tools.JsonUnMarshal([]byte(setting.Value), &adminIds)
				if err == nil {
					for _, id := range adminIds {
						administrators[id] = true
					}
					ret.GroupManagement.AdminCount = int32(len(adminIds))
				}
			}
		} else if setting.Key == string(commonservices.AttItemKey_GrpVerifyType) {
			verifyType := tools.ToInt(setting.Value)
			ret.GroupManagement.GroupVerifyType = int32(verifyType)
		} else if setting.Key == string(commonservices.AttItemKey_HideGrpMsg) {
			hidGrpMsg := tools.ToInt(setting.Value)
			var visible int32 = 0
			if hidGrpMsg > 0 {
				visible = 0
			} else {
				visible = 1
			}
			ret.GroupManagement.GroupHisMsgVisible = visible
		}
	}
	// my role
	myRole := pbobjs.GrpMemberRole_GrpMember // 0: 群成员；1:群主；2:群管理员；3:非群成员；
	if requestId == creator {
		myRole = pbobjs.GrpMemberRole_GrpCreator
	} else if _, exist := administrators[requestId]; exist {
		myRole = pbobjs.GrpMemberRole_GrpAdmin
	}
	ret.MyRole = myRole
	//owner
	if creator != "" {
		ownerUser := commonservices.GetTargetDisplayUserInfo(ctx, creator)
		if ownerUser != nil {
			ret.Owner = &pbobjs.GroupMemberInfo{
				UserId:     ownerUser.UserId,
				Nickname:   ownerUser.Nickname,
				Avatar:     ownerUser.UserPortrait,
				MemberType: ownerUser.UserType,
			}
		}
	}
	//top members
	code, topMembers := QueryGrpMembers(ctx, &pbobjs.QryGroupMembersReq{
		GroupId: groupId,
		Limit:   20,
	})
	if code == errs.IMErrorCode_SUCCESS && topMembers != nil {
		for _, member := range topMembers.Items {
			member.Role = pbobjs.GrpMemberRole_GrpMember
			if member.UserId == creator {
				member.Role = pbobjs.GrpMemberRole_GrpCreator
			} else if _, exist := administrators[member.UserId]; exist {
				member.Role = pbobjs.GrpMemberRole_GrpAdmin
			}
			ret.Members = append(ret.Members, member)
		}
	}
	//qry group member exts/settings
	code, respObj, err = bases.SyncRpcCall(ctx, "qry_grp_member_settings", groupId, &pbobjs.QryGrpMemberSettingsReq{
		MemberId: requestId,
	}, func() proto.Message {
		return &pbobjs.QryGrpMemberSettingsResp{}
	})
	if err == nil && code == errs.IMErrorCode_SUCCESS && respObj != nil {
		memberSettings := respObj.(*pbobjs.QryGrpMemberSettingsResp)
		if displayName, exist := memberSettings.MemberExts[string(commonservices.AttItemKey_GrpDisplayName)]; exist {
			ret.GrpDisplayName = displayName
		}
	}

	return errs.IMErrorCode_SUCCESS, ret
}

func CheckGroupMembers(ctx context.Context, req *models.CheckGroupMembersReq) (errs.IMErrorCode, *models.CheckGroupMembersResp) {
	ret := &models.CheckGroupMembersResp{
		GroupId:        req.GroupId,
		MemberExistMap: map[string]bool{},
	}
	for _, memberId := range req.MemberIds {
		ret.MemberExistMap[memberId] = false
	}
	code, respObj, err := bases.SyncRpcCall(ctx, "g_check_members", req.GroupId, &pbobjs.CheckGroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: req.MemberIds,
	}, func() proto.Message {
		return &pbobjs.CheckGroupMembersResp{}
	})
	if err == nil && code == errs.IMErrorCode_SUCCESS && respObj != nil {
		checkResp, ok := respObj.(*pbobjs.CheckGroupMembersResp)
		if ok {
			for memberId := range checkResp.MemberIdMap {
				ret.MemberExistMap[memberId] = true
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func CreateGroup(ctx context.Context, req *pbobjs.GroupMembersReq) (errs.IMErrorCode, *pbobjs.GroupInfo) {
	grpId := tools.GenerateUUIDShort11()
	requestId := bases.GetRequesterIdFromCtx(ctx)
	memberIds := []string{requestId}
	for _, memberId := range req.MemberIds {
		if memberId != requestId {
			memberIds = append(memberIds, memberId)
		}
	}
	settings := []*pbobjs.KvItem{}
	settings = append(settings, &pbobjs.KvItem{
		Key:   string(commonservices.AttItemKey_GrpCreator),
		Value: requestId,
	})
	code, _, err := bases.SyncRpcCall(ctx, "g_add_members", grpId, &pbobjs.GroupMembersReq{
		GroupId:       grpId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
		MemberIds:     memberIds,
		Settings:      settings,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	//send notify msg
	targetUsers := []*pbobjs.UserObj{}
	for _, memberId := range req.MemberIds {
		targetUsers = append(targetUsers, GetUser(ctx, memberId))
	}
	notify := &models.GroupNotify{
		Operator: GetUser(ctx, requestId),
		Members:  targetUsers,
		Type:     models.GroupNotifyType_AddMember,
	}
	SendGrpNotify(ctx, grpId, notify)
	return errs.IMErrorCode_SUCCESS, &pbobjs.GroupInfo{
		GroupId:       grpId,
		GroupName:     req.GroupName,
		GroupPortrait: req.GroupPortrait,
	}
}

func UpdateGroup(ctx context.Context, req *pbobjs.GroupInfo) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := bases.SyncRpcCall(ctx, "upd_group_info", req.GroupId, req, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	SendGrpNotify(ctx, req.GroupId, &models.GroupNotify{
		Operator: GetUser(ctx, requestId),
		Name:     req.GroupName,
		Type:     models.GroupNotifyType_Rename,
	})
	return errs.IMErrorCode_SUCCESS
}

func DissolveGroup(ctx context.Context, groupId string) errs.IMErrorCode {
	code, _, err := bases.SyncRpcCall(ctx, "g_dissolve", groupId, &pbobjs.GroupMembersReq{
		GroupId: groupId,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	return errs.IMErrorCode_SUCCESS
}

func QuitGroup(ctx context.Context, groupId string) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := bases.SyncRpcCall(ctx, "g_del_members", groupId, &pbobjs.GroupMembersReq{
		GroupId:   groupId,
		MemberIds: []string{requestId},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	SendGrpNotify(ctx, groupId, &models.GroupNotify{
		Operator: GetUser(ctx, requestId),
		Members: []*pbobjs.UserObj{
			GetUser(ctx, requestId),
		},
		Type: models.GroupNotifyType_RemoveMember,
	})
	return errs.IMErrorCode_SUCCESS
}

func AddGrpMembers(ctx context.Context, grpMembers *pbobjs.GroupMembersReq) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := bases.SyncRpcCall(ctx, "g_add_members", grpMembers.GroupId, &pbobjs.GroupMembersReq{
		GroupId:   grpMembers.GroupId,
		MemberIds: grpMembers.MemberIds,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	//send notify msg
	targetUsers := []*pbobjs.UserObj{}
	for _, memberId := range grpMembers.MemberIds {
		targetUsers = append(targetUsers, GetUser(ctx, memberId))
	}
	notify := &models.GroupNotify{
		Operator: GetUser(ctx, requestId),
		Members:  targetUsers,
		Type:     models.GroupNotifyType_AddMember,
	}
	//send notify msg
	SendGrpNotify(ctx, grpMembers.GroupId, notify)
	return errs.IMErrorCode_SUCCESS
}

func GrpInviteMembers(ctx context.Context, req *pbobjs.GroupInviteReq) (errs.IMErrorCode, *pbobjs.GroupInviteResp) {
	appkey := bases.GetRequesterIdFromCtx(ctx)
	requesterId := bases.GetRequesterIdFromCtx(ctx)
	//TODO check operator
	results := &pbobjs.GroupInviteResp{
		Results: make(map[string]pbobjs.GrpInviteResultReason),
	}
	//TODO check grp member exist
	//check user's setting
	directAddMemberIds := []string{}
	for _, memberId := range req.MemberIds {
		reason := pbobjs.GrpInviteResultReason_InviteSucc
		mUserInfo := commonservices.GetTargetUserInfo(ctx, memberId)
		mUserSetting := GetUserSettings(mUserInfo)
		if mUserSetting.GrpVerifyType == pbobjs.GrpVerifyType_DeclineGroup {
			reason = pbobjs.GrpInviteResultReason_InviteDecline
		} else if mUserSetting.GrpVerifyType == pbobjs.GrpVerifyType_NeedGrpVerify {
			storage := storages.NewGrpApplicationStorage()
			storage.InviteUpsert(storeModels.GrpApplication{
				GroupId:     req.GroupId,
				ApplyType:   storeModels.GrpApplicationType_Invite,
				RecipientId: memberId,
				InviterId:   requesterId,
				ApplyTime:   time.Now().UnixMilli(),
				Status:      storeModels.GrpApplicationStatus_Invite,
				AppKey:      appkey,
			})
			reason = pbobjs.GrpInviteResultReason_InviteSendOut
		} else if mUserSetting.GrpVerifyType == pbobjs.GrpVerifyType_NoNeedGrpVerify {
			directAddMemberIds = append(directAddMemberIds, memberId)
			reason = pbobjs.GrpInviteResultReason_InviteSucc
		}
		results.Results[memberId] = reason
	}
	if len(directAddMemberIds) > 0 {
		code, _, err := bases.SyncRpcCall(ctx, "g_add_members", req.GroupId, &pbobjs.GroupMembersReq{
			GroupId:   req.GroupId,
			MemberIds: directAddMemberIds,
		}, nil)
		if err != nil || code != errs.IMErrorCode_SUCCESS {
			return code, nil
		}
		//send notify msg
		targetUsers := []*pbobjs.UserObj{}
		for _, memberId := range directAddMemberIds {
			targetUsers = append(targetUsers, GetUser(ctx, memberId))
		}
		notify := &models.GroupNotify{
			Operator: GetUser(ctx, requesterId),
			Members:  targetUsers,
			Type:     models.GroupNotifyType_AddMember,
		}
		SendGrpNotify(ctx, req.GroupId, notify)
	}
	return errs.IMErrorCode_SUCCESS, results
}

func GrpJoinApply(ctx context.Context, req *pbobjs.GroupInviteReq) errs.IMErrorCode {
	userId := bases.GetRequesterIdFromCtx(ctx)
	groupId := req.GroupId
	//check grp member exists
	code, respObj, err := bases.SyncRpcCall(ctx, "g_check_members", groupId, &pbobjs.CheckGroupMembersReq{
		GroupId:   groupId,
		MemberIds: []string{userId},
	}, func() proto.Message {
		return &pbobjs.CheckGroupMembersResp{}
	})
	if err == nil && code == errs.IMErrorCode_SUCCESS && respObj != nil {
		checkResp, ok := respObj.(*pbobjs.CheckGroupMembersResp)
		if ok {
			if _, exist := checkResp.MemberIdMap[userId]; exist {
				return errs.IMErrorCode_APP_GROUP_MEMBEREXISTED
			}
		}
	}
	//add group
	code, _, err = bases.SyncRpcCall(ctx, "g_add_members", groupId, &pbobjs.GroupMembersReq{
		GroupId:   groupId,
		MemberIds: []string{userId},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	//send notify msg
	notify := &models.GroupNotify{
		Operator: GetUser(ctx, userId),
		Type:     models.GroupNotifyType_Join,
	}
	SendGrpNotify(ctx, groupId, notify)
	return errs.IMErrorCode_SUCCESS
}

func DelGrpMembers(ctx context.Context, req *pbobjs.GroupMembersReq) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := bases.SyncRpcCall(ctx, "g_del_members", req.GroupId, &pbobjs.GroupMembersReq{
		GroupId:   req.GroupId,
		MemberIds: req.MemberIds,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	//send notify msg
	targetUsers := []*pbobjs.UserObj{}
	for _, memberId := range req.MemberIds {
		targetUsers = append(targetUsers, GetUser(ctx, memberId))
	}
	SendGrpNotify(ctx, req.GroupId, &models.GroupNotify{
		Operator: GetUser(ctx, requestId),
		Members:  targetUsers,
		Type:     models.GroupNotifyType_RemoveMember,
	})
	return errs.IMErrorCode_SUCCESS
}

func QueryGrpMembers(ctx context.Context, req *pbobjs.QryGroupMembersReq) (errs.IMErrorCode, *pbobjs.GroupMemberInfos) {
	code, respObj, err := bases.SyncRpcCall(ctx, "g_qry_members", req.GroupId, req, func() proto.Message {
		return &pbobjs.GroupMembersResp{}
	})
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code, nil
	}
	members := respObj.(*pbobjs.GroupMembersResp)
	ret := &pbobjs.GroupMemberInfos{
		Items:  []*pbobjs.GroupMemberInfo{},
		Offset: members.Offset,
	}
	memberIds := []string{}
	for _, member := range members.Items {
		memberIds = append(memberIds, member.MemberId)
		ret.Items = append(ret.Items, &pbobjs.GroupMemberInfo{
			UserId:     member.MemberId,
			MemberType: member.MemberType,
		})
	}
	userMap := commonservices.GetTargetDisplayUserInfosMap(ctx, memberIds)
	for _, member := range ret.Items {
		userInfo, ok := userMap[member.UserId]
		if ok && userInfo != nil {
			member.Nickname = userInfo.Nickname
			member.Avatar = userInfo.UserPortrait
			member.MemberType = userInfo.UserType
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func SetGrpAnnouncement(ctx context.Context, req *pbobjs.GrpAnnouncement) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := bases.SyncRpcCall(ctx, "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId: req.GroupId,
		Settings: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_GrpAnnouncement),
				Value: req.Content,
			},
		},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	if req.Content != "" {
		//send announce msg
		flag := msgdefines.SetStoreMsg(0)
		flag = msgdefines.SetCountMsg(flag)
		txtMsg := &msgdefines.TextMsg{
			Content: req.Content,
		}
		commonservices.AsyncGroupMsgOverUpstream(ctx, requestId, req.GroupId, &pbobjs.UpMsg{
			MsgType:    msgdefines.InnerMsgType_Text,
			MsgContent: []byte(tools.ToJson(txtMsg)),
			Flags:      flag,
			MentionInfo: &pbobjs.MentionInfo{
				MentionType: pbobjs.MentionType_All,
			},
		})
	}
	return errs.IMErrorCode_SUCCESS
}

func GetGrpAnnouncement(ctx context.Context, groupId string) (errs.IMErrorCode, *pbobjs.GrpAnnouncement) {
	grpInfo := commonservices.GetGroupInfoFromCache(ctx, groupId)
	ret := &pbobjs.GrpAnnouncement{
		GroupId: groupId,
	}
	if grpInfo != nil {
		for _, kv := range grpInfo.Settings {
			if kv.Key == string(commonservices.AttItemKey_GrpAnnouncement) {
				ret.Content = kv.Value
			}
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func ChgGroupOwner(ctx context.Context, req *pbobjs.GroupOwnerChgReq) errs.IMErrorCode {
	//TODO check right
	code, _, err := bases.SyncRpcCall(ctx, "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId: req.GroupId,
		Settings: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_GrpCreator),
				Value: req.OwnerId,
			},
		},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	//send notify
	requestId := bases.GetRequesterIdFromCtx(ctx)
	notify := &models.GroupNotify{
		Operator: GetUser(ctx, requestId),
		Members: []*pbobjs.UserObj{
			GetUser(ctx, req.OwnerId),
		},
		Type: models.GroupNotifyType_ChgOwner,
	}
	SendGrpNotify(ctx, req.GroupId, notify)
	return errs.IMErrorCode_SUCCESS
}

func SetGroupMute(ctx context.Context, req *pbobjs.SetGroupMuteReq) errs.IMErrorCode {
	//TODO check right
	code, _, err := bases.SyncRpcCall(ctx, "group_mute", req.GroupId, &pbobjs.GroupMuteReq{
		GroupId: req.GroupId,
		IsMute:  req.IsMute,
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupVerifyType(ctx context.Context, req *pbobjs.SetGroupVerifyTypeReq) errs.IMErrorCode {
	//TODO check right
	code, _, err := bases.SyncRpcCall(ctx, "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId: req.GroupId,
		Settings: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_GrpVerifyType),
				Value: tools.Int642String(int64(req.VerifyType)),
			},
		},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	return errs.IMErrorCode_SUCCESS
}

func SetGroupHisMsgVisible(ctx context.Context, req *pbobjs.SetGroupHisMsgVisibleReq) errs.IMErrorCode {
	//TODO check right
	visible := req.GroupHisMsgVisible
	hideGrpMsg := "1"
	if visible > 0 {
		hideGrpMsg = "0"
	} else {
		hideGrpMsg = "1"
	}
	code, _, err := bases.SyncRpcCall(ctx, "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId: req.GroupId,
		Settings: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_HideGrpMsg),
				Value: hideGrpMsg,
			},
		},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	return errs.IMErrorCode_SUCCESS
}

func AddGroupAdministrators(ctx context.Context, req *pbobjs.GroupAdministratorsReq) errs.IMErrorCode {
	grpInfo := commonservices.GetGroupInfoFromRpc(ctx, req.GroupId)
	adminIds := []string{}
	if grpInfo != nil {
		for _, setting := range grpInfo.Settings {
			if setting.Key == string(commonservices.AttItemKey_GrpAdministrators) {
				tools.JsonUnMarshal([]byte(setting.Value), &adminIds)
				break
			}
		}
	}
	if len(adminIds)+len(req.AdminIds) > 10 {
		return errs.IMErrorCode_APP_DEFAULT
	}
	adminIds = append(adminIds, req.AdminIds...)
	code, _, err := bases.SyncRpcCall(ctx, "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId: req.GroupId,
		Settings: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_GrpAdministrators),
				Value: tools.ToJson(adminIds),
			},
		},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	return errs.IMErrorCode_SUCCESS
}

func DelGroupAdministrators(ctx context.Context, req *pbobjs.GroupAdministratorsReq) errs.IMErrorCode {
	grpInfo := commonservices.GetGroupInfoFromRpc(ctx, req.GroupId)
	adminIds := []string{}
	if grpInfo != nil {
		for _, setting := range grpInfo.Settings {
			if setting.Key == string(commonservices.AttItemKey_GrpAdministrators) {
				tools.JsonUnMarshal([]byte(setting.Value), &adminIds)
				break
			}
		}
	}
	needDelMap := map[string]int{}
	for _, id := range req.AdminIds {
		needDelMap[id] = 1
	}
	newAdminIds := []string{}
	for _, id := range adminIds {
		if _, exist := needDelMap[id]; !exist {
			newAdminIds = append(newAdminIds, id)
		}
	}
	code, _, err := bases.SyncRpcCall(ctx, "upd_group_info", req.GroupId, &pbobjs.GroupInfo{
		GroupId: req.GroupId,
		Settings: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_GrpAdministrators),
				Value: tools.ToJson(newAdminIds),
			},
		},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	return errs.IMErrorCode_SUCCESS
}

func QryGroupAdministrators(ctx context.Context, groupId string) (errs.IMErrorCode, *pbobjs.GroupAdministratorsResp) {
	ret := &pbobjs.GroupAdministratorsResp{
		GroupId: groupId,
		Items:   []*pbobjs.GroupMemberInfo{},
	}
	grpInfo := commonservices.GetGroupInfoFromRpc(ctx, groupId)
	if grpInfo != nil {
		adminIds := []string{}
		for _, setting := range grpInfo.Settings {
			if setting.Key == string(commonservices.AttItemKey_GrpAdministrators) {
				tools.JsonUnMarshal([]byte(setting.Value), &adminIds)
				break
			}
		}
		userMap := commonservices.GetTargetDisplayUserInfosMap(ctx, adminIds)
		for _, userId := range adminIds {
			grpMember := &pbobjs.GroupMemberInfo{
				UserId: userId,
			}
			if userInfo, exist := userMap[userId]; exist {
				grpMember.Nickname = userInfo.Nickname
				grpMember.Avatar = userInfo.UserPortrait
			}
			ret.Items = append(ret.Items, grpMember)
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func SetGrpDisplayName(ctx context.Context, req *pbobjs.SetGroupDisplayNameReq) errs.IMErrorCode {
	requestId := bases.GetRequesterIdFromCtx(ctx)
	code, _, err := bases.SyncRpcCall(ctx, "set_grp_member_setting", req.GroupId, &pbobjs.GroupMember{
		MemberId: requestId,
		ExtFields: []*pbobjs.KvItem{
			{
				Key:   string(commonservices.AttItemKey_GrpDisplayName),
				Value: req.GrpDisplayName,
			},
		},
	}, nil)
	if err != nil || code != errs.IMErrorCode_SUCCESS {
		return code
	}
	return errs.IMErrorCode_SUCCESS
}

func QryMyGrpApplications(ctx context.Context, req *pbobjs.QryGrpApplicationsReq) (errs.IMErrorCode, *pbobjs.QryGrpApplicationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewGrpApplicationStorage()
	ret := &pbobjs.QryGrpApplicationsResp{
		Items: []*pbobjs.GrpApplicationItem{},
	}
	applications, err := storage.QueryMyGrpApplications(appkey, userId, req.StartTime, int64(req.Count), req.Order > 0)
	if err == nil {
		for _, application := range applications {
			ret.Items = append(ret.Items, &pbobjs.GrpApplicationItem{
				GrpInfo: &pbobjs.GrpInfo{
					GroupId: application.GroupId,
				},
				ApplyType: int32(application.ApplyType),
				Operator: &pbobjs.UserObj{
					UserId: application.OperatorId,
				},
				ApplyTime: application.ApplyTime,
				Status:    int32(application.Status),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryMyPendingGrpInvitations(ctx context.Context, req *pbobjs.QryGrpApplicationsReq) (errs.IMErrorCode, *pbobjs.QryGrpApplicationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	userId := bases.GetRequesterIdFromCtx(ctx)
	storage := storages.NewGrpApplicationStorage()
	ret := &pbobjs.QryGrpApplicationsResp{
		Items: []*pbobjs.GrpApplicationItem{},
	}
	applications, err := storage.QueryMyPendingGrpInvitations(appkey, userId, req.StartTime, int64(req.Count), req.Order > 0)
	if err == nil {
		for _, application := range applications {
			ret.Items = append(ret.Items, &pbobjs.GrpApplicationItem{
				GrpInfo: &pbobjs.GrpInfo{
					GroupId: application.GroupId,
				},
				ApplyType: int32(application.ApplyType),
				Inviter: &pbobjs.UserObj{
					UserId: application.InviterId,
				},
				ApplyTime: application.ApplyTime,
				Status:    int32(application.Status),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryGrpInvitations(ctx context.Context, req *pbobjs.QryGrpApplicationsReq) (errs.IMErrorCode, *pbobjs.QryGrpApplicationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewGrpApplicationStorage()
	ret := &pbobjs.QryGrpApplicationsResp{
		Items: []*pbobjs.GrpApplicationItem{},
	}
	applications, err := storage.QueryGrpInvitations(appkey, req.GroupId, req.StartTime, int64(req.Count), req.Order > 0)
	if err == nil {
		for _, application := range applications {
			ret.Items = append(ret.Items, &pbobjs.GrpApplicationItem{
				GrpInfo: &pbobjs.GrpInfo{
					GroupId: application.GroupId,
				},
				ApplyType: int32(application.ApplyType),
				Recipient: &pbobjs.UserObj{
					UserId: application.RecipientId,
				},
				Inviter: &pbobjs.UserObj{
					UserId: application.InviterId,
				},
				ApplyTime: application.ApplyTime,
				Status:    int32(application.Status),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}

func QryGrpPendingApplications(ctx context.Context, req *pbobjs.QryGrpApplicationsReq) (errs.IMErrorCode, *pbobjs.QryGrpApplicationsResp) {
	appkey := bases.GetAppKeyFromCtx(ctx)
	storage := storages.NewGrpApplicationStorage()
	ret := &pbobjs.QryGrpApplicationsResp{
		Items: []*pbobjs.GrpApplicationItem{},
	}
	applications, err := storage.QueryGrpPendingApplications(appkey, req.GroupId, req.StartTime, int64(req.Count), req.Order > 0)
	if err == nil {
		for _, application := range applications {
			ret.Items = append(ret.Items, &pbobjs.GrpApplicationItem{
				GrpInfo: &pbobjs.GrpInfo{
					GroupId: application.GroupId,
				},
				ApplyType: int32(application.ApplyType),
				Sponsor: &pbobjs.UserObj{
					UserId: application.SponsorId,
				},
				Operator: &pbobjs.UserObj{
					UserId: application.OperatorId,
				},
				ApplyTime: application.ApplyTime,
				Status:    int32(application.Status),
			})
		}
	}
	return errs.IMErrorCode_SUCCESS, ret
}
