package dbs

import (
	"fmt"
	"im-server/commons/dbcommons"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/rtcroom/storages/models"
	"time"
)

type RtcRoomMemberDao struct {
	ID             int64  `gorm:"primary_key"`
	RoomId         string `gorm:"room_id"`
	MemberId       string `gorm:"member_id"`
	DeviceId       string `gorm:"device_id"`
	RtcState       int    `gorm:"rtc_state"`
	InviterId      string `gorm:"inviter_id"`
	LatestPingTime int64  `gorm:"latest_ping_time"`
	CameraEnable   int32  `gorm:"camera_enable"`
	MicEnable      int32  `gorm:"mic_enable"`
	CallTime       int64  `gorm:"call_time"`
	ConnectTime    int64  `gorm:"connect_time"`
	HangupTime     int64  `gorm:"hangup_time"`
	AppKey         string `gorm:"app_key"`
}

func (member *RtcRoomMemberDao) TableName() string {
	return "rtcmembers"
}

func (member *RtcRoomMemberDao) Upsert(item models.RtcRoomMember) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_key,room_id,member_id,device_id,rtc_state,inviter_id,camera_enable,mic_enable,call_time,connect_time,hangup_time,latest_ping_time)VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE device_id=?,rtc_state=?,camera_enable=?,mic_enable=?,call_time=?,connect_time=?,hangup_time=?,latest_ping_time=?", member.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppKey, item.RoomId, item.MemberId, item.DeviceId, int(item.RtcState), item.InviterId, item.CameraEnable, item.MicEnable, item.CallTime, item.ConnectTime, item.HangupTime, item.LatestPingTime,
		item.DeviceId, int(item.RtcState), item.CameraEnable, item.MicEnable, item.CallTime, item.ConnectTime, item.HangupTime, item.LatestPingTime).Error
}

func (member *RtcRoomMemberDao) Insert(item models.RtcRoomMember) (int64, error) {
	sql := fmt.Sprintf("INSERT IGNORE INTO %s (app_key,room_id,member_id,device_id,rtc_state,inviter_id,call_time,connect_time,hangup_time,latest_ping_time)VALUES(?,?,?,?,?,?,?,?,?,?)", member.TableName())
	result := dbcommons.GetDb().Exec(sql, item.AppKey, item.RoomId, item.MemberId, item.DeviceId, item.RtcState, item.InviterId, item.CallTime, item.ConnectTime, item.HangupTime, item.LatestPingTime)
	return result.RowsAffected, result.Error
}

func (member *RtcRoomMemberDao) Find(appkey, roomId, memberId string) (*models.RtcRoomMember, error) {
	var item RtcRoomMemberDao
	err := dbcommons.GetDb().Where("app_key=? and room_id=? and member_id=?", appkey, roomId, memberId).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &models.RtcRoomMember{
		RoomId:    item.RoomId,
		MemberId:  item.MemberId,
		DeviceId:  item.DeviceId,
		RtcState:  pbobjs.RtcState(item.RtcState),
		InviterId: item.InviterId,
		AppKey:    item.AppKey,
	}, nil
}

func (member *RtcRoomMemberDao) UpdateState(appkey, roomId, memberId string, state pbobjs.RtcState, deviceId string) error {
	upd := map[string]interface{}{}
	if state != pbobjs.RtcState_RtcStateDefault {
		upd["rtc_state"] = state
	}
	if deviceId != "" {
		upd["device_id"] = deviceId
	}
	upd["latest_ping_time"] = time.Now().UnixMilli()
	return dbcommons.GetDb().Model(member).Where("app_key=? and room_id=? and member_id=?", appkey, roomId, memberId).Update(upd).Error
}

func (member *RtcRoomMemberDao) RefreshPingTime(appkey, roomId, memberId string) error {
	upd := map[string]interface{}{}
	upd["latest_ping_time"] = time.Now().UnixMilli()
	return dbcommons.GetDb().Model(member).Where("app_key=? and room_id=? and member_id=?", appkey, roomId, memberId).Update(upd).Error
}

func (member *RtcRoomMemberDao) Delete(appkey, roomId, memberId string) error {
	return dbcommons.GetDb().Where("app_key=? and room_id=? and member_id=?", appkey, roomId, memberId).Delete(&RtcRoomMemberDao{}).Error
}

func (member *RtcRoomMemberDao) DeleteByRoomId(appkey, roomId string) error {
	return dbcommons.GetDb().Where("app_key=? and room_id=?", appkey, roomId).Delete(&RtcRoomMemberDao{}).Error
}

func (member *RtcRoomMemberDao) QueryMembers(appkey, roomId string, startId, limit int64) ([]*models.RtcRoomMember, error) {
	var items []*RtcRoomMemberDao
	err := dbcommons.GetDb().Where("app_key=? and room_id=? and id>?", appkey, roomId, startId).Order("id asc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.RtcRoomMember{}
	for _, item := range items {
		ret = append(ret, &models.RtcRoomMember{
			ID:           item.ID,
			RoomId:       item.RoomId,
			MemberId:     item.MemberId,
			DeviceId:     item.DeviceId,
			RtcState:     pbobjs.RtcState(item.RtcState),
			InviterId:    item.InviterId,
			CameraEnable: item.CameraEnable,
			MicEnable:    item.MicEnable,
			CallTime:     item.CallTime,
			ConnectTime:  item.ConnectTime,
			HangupTime:   item.HangupTime,
			AppKey:       item.AppKey,
		})
	}
	return ret, nil
}

func (member *RtcRoomMemberDao) QueryRoomsByMember(appkey, memberId string, limit int64) ([]*models.RtcRoomMember, error) {
	var items []*RtcRoomMemberDao
	err := dbcommons.GetDb().Where("app_key=? and member_id=?", appkey, memberId).Order("id desc").Limit(limit).Find(&items).Error
	if err != nil {
		return nil, err
	}
	ret := []*models.RtcRoomMember{}
	for _, item := range items {
		ret = append(ret, &models.RtcRoomMember{
			ID:             item.ID,
			RoomId:         item.RoomId,
			MemberId:       item.MemberId,
			DeviceId:       item.DeviceId,
			RtcState:       pbobjs.RtcState(item.RtcState),
			InviterId:      item.InviterId,
			LatestPingTime: item.LatestPingTime,
			AppKey:         item.AppKey,
		})
	}
	return ret, nil
}
