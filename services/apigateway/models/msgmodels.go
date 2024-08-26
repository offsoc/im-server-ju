package models

type SendMsgReq struct {
	SenderId       string       `json:"sender_id"`
	ReceiverId     string       `json:"receiver_id"`
	TargetId       string       `json:"target_id"`
	TargetIds      []string     `json:"target_ids"`
	MsgType        string       `json:"msg_type"`
	MsgContent     string       `json:"msg_content"`
	IsStorage      *bool        `json:"is_storage"`
	IsCount        *bool        `json:"is_count"`
	IsNotifySender *bool        `json:"is_notify_sender"`
	IsState        *bool        `json:"is_state"`
	IsCmd          *bool        `json:"is_cmd"`
	MentionInfo    *MentionInfo `json:"mention_info"`
	ReferMsg       *ReferMsg    `json:"refer_msg"`
}

type MentionInfo struct {
	MentionType   string      `json:"mention_type"`
	TargetUsers   []*UserInfo `json:"target_users"`
	TargetUserIds []string    `json:"target_user_ids"`
}

type ReferMsg struct {
	MsgId       string `json:"msg_id"`
	SenderId    string `json:"sender_id"`
	TargetId    string `json:"target_id"`
	ChannelType int    `json:"channel_type"`
	MsgType     string `json:"msg_type"`
	MsgTime     int64  `json:"msg_time"`
	MsgContent  string `json:"msg_content"`
}

type SendGrpCastMsgReq struct {
	SenderId      string          `json:"sender_id"`
	TargetId      string          `json:"target_id"`
	MsgType       string          `json:"msg_type"`
	MsgContent    string          `json:"msg_content"`
	TargetConvers []*Conversation `json:"target_convers"`
}

type SendBrdCastMsgReq struct {
	SenderId   string `json:"sender_id"`
	MsgType    string `json:"msg_type"`
	MsgContent string `json:"msg_content"`
	IsStorage  *bool  `json:"is_storage"`
}

/*
MentionType mentionType = 1;

	repeated UserInfo targetUsers = 2;
*/
type SendMsgResp struct {
	MsgId string `json:"msg_id"`
}

type HisMsgs struct {
	Msgs       []*HisMsg `json:"msgs"`
	IsFinished bool      `json:"is_finished"`
}
type HisMsg struct {
	SenderId    string `json:"sender_id"`
	ReceiverId  string `json:"receiver_id"`
	ChannelType int32  `json:"channel_type"`
	MsgId       string `json:"msg_id"`
	MsgTime     int64  `json:"msg_time"`
	MsgType     string `json:"msg_type"`
	MsgContent  string `json:"msg_content"`
}

type CleanHisMsgsReq struct {
	FromId          string `json:"from_id"`
	TargetId        string `json:"target_id"`
	ChannelType     int32  `json:"channel_type"`
	CleanTime       int64  `json:"clean_time"`
	CleanTimeOffset int64  `json:"clean_time_offset"`
	CleanScope      int    `json:"clean_scope"`
	SenderId        string `json:"sender_id,omitempty"`
}

type RecallHisMsgsReq struct {
	FromId      string            `json:"from_id"`
	TargetId    string            `json:"target_id"`
	ChannelType int32             `json:"channel_type"`
	MsgId       string            `json:"msg_id"`
	MsgTime     int64             `json:"msg_time"`
	Exts        map[string]string `json:"exts"`
}

type DelHisMsgsReq struct {
	FromId      string       `json:"from_id"`
	TargetId    string       `json:"target_id"`
	ChannelType int32        `json:"channel_type"`
	DelScope    int          `json:"del_scope"`
	Msgs        []*SimpleMsg `json:"msgs"`
}
type SimpleMsg struct {
	MsgId        string `json:"msg_id"`
	MsgTime      int64  `json:"msg_time"`
	MsgReadIndex int64  `json:"msg_read_index"`
}
