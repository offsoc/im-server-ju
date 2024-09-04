package main

import (
	"fmt"
	"im-server/commons/pbdefines/pbobjs"
	"im-server/services/commonservices"
	"im-server/simulator/wsclients"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type (
	groupMessageData struct {
		msgId       string
		sendTime    time.Time
		receiveList []receiveData
	}
	receiveData struct {
		Time     time.Time
		TargetId string
	}
)

func groupSend() {
	var sendClients = make(map[int]*wsclients.WsImClient, 1000)
	var clientLocker sync.Mutex

	var connectWg sync.WaitGroup
	var sendWg sync.WaitGroup

	var (
		msgMap     = make(map[string]*groupMessageData, 10000)
		msgMapLock sync.Mutex
	)

	var (
		groupMemberNum = 1000
		sendClientNum  = 10
		turnMsgCount   = 200
	)
	//添加群组
	groupId := "benchmark_group:" + strconv.Itoa(groupMemberNum)
	dissolveGroup(groupId)
	createGroup(groupId, groupMemberNum)

	for i := 0; i < groupMemberNum; i++ {
		connectWg.Add(1)

		go func(i int) {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			client := wsclients.NewWsImClient(wsUrl, appKey, generateUserTokenStr(userPrefix+strconv.Itoa(i)), func(msg *pbobjs.DownMsg) {
				sendTime := bytesToTime(msg.MsgContent)
				msgMapLock.Lock()
				if _, ok := msgMap[msg.MsgId]; ok {
					data, ok := msgMap[msg.MsgId]
					if !ok || data == nil {
						data = &groupMessageData{
							msgId:    msg.MsgId,
							sendTime: sendTime,
						}
					}
					data.receiveList = append(data.receiveList, receiveData{
						Time:     time.Now(),
						TargetId: userPrefix + strconv.Itoa(i),
					})
				}
				msgMapLock.Unlock()

				//fmt.Printf("msgId:%s, senderId:%s, receiverId:%s, msgTime:%v, msgContent:%s\n", msg.MsgId, msg.SenderId, msg.TargetId, msg.MsgTime, msg.MsgContent)

			}, func(msg *pbobjs.StreamDownMsg) {
				fmt.Println("stream down msg:", msg)
			}, nil)
			code, _ := client.Connect("", "")
			if code != 0 {
				fmt.Printf("connect code: %d\n", code)
			}
			clientLocker.Lock()
			if _, ok := sendClients[i]; ok {
				panic("client already exists")
			}
			sendClients[i] = client
			clientLocker.Unlock()
			connectWg.Done()
		}(i)
	}
	connectWg.Wait()
	defer func() {
		for _, client := range sendClients {
			client.Disconnect()
		}
	}()
	fmt.Printf("连接创建完成，连接数 %d\n", len(sendClients))

	sendStart := time.Now()

	fmt.Printf("开始发送消息, 发送客户端数量 %d, 发送消息数量 %d\n", sendClientNum, turnMsgCount)
	for i := 0; i < sendClientNum; i++ {
		sendWg.Add(1)
		go func(i int) {
			client := sendClients[i]
			flag := commonservices.SetStoreMsg(0)
			flag = commonservices.SetCountMsg(flag)

			for j := 0; j < turnMsgCount; j++ {
				upMsg := pbobjs.UpMsg{
					MsgType:    "txtMsg",
					MsgContent: timeToBytes(time.Now()),
					Flags:      commonservices.SetStoreMsg(0),
					MentionInfo: &pbobjs.MentionInfo{
						MentionType: pbobjs.MentionType_All,
					},
				}
				code, ack := client.SendGroupMsg(groupId, &upMsg)
				if code != 0 {
					fmt.Printf("send upMsg failed, code: %d\n", code)
				}

				msgMapLock.Lock()
				msgMap[ack.MsgId] = nil
				msgMapLock.Unlock()
			}
			sendWg.Done()
		}(i)
	}
	WaitTimeout(&sendWg, 10*time.Second)
	fmt.Printf("发送消息数量 %d, time used:%v\n", len(msgMap), time.Now().Sub(sendStart))

	s := time.Now()
	ticker := time.NewTicker(time.Second)

	prevTotal := 0
	for {
		select {
		case t := <-ticker.C:
			if t.Sub(s).Seconds() > 10 {
				return
			}

			msgMapLock.Lock()
			total, maxDelay, minDelay, avgDelay := statisticsGroupMsgMap(msgMap)
			if total > prevTotal {
				fmt.Printf("收到消息数量 %d, 平均延迟 %v, 最大延迟 %v, 最小延迟 %v\n", total, avgDelay, maxDelay, minDelay)
				prevTotal = total
			}
			if total >= len(msgMap) {
				return
			}
			msgMapLock.Unlock()
		}
	}

}

func statisticsGroupMsgMap(msgMap map[string]*groupMessageData) (total int, maxDelay time.Duration, minDelay time.Duration, avgDelay time.Duration) {
	var (
		totalReceiveCount int
		totalDelay        time.Duration
	)
	for _, msg := range msgMap {
		if msg == nil {
			continue
		}
		total++
		totalReceiveCount += len(msg.receiveList)

		for _, data := range msg.receiveList {
			duration := data.Time.Sub(msg.sendTime)

			totalDelay += duration
			if duration > maxDelay {
				maxDelay = duration
			}
			if duration < minDelay {
				minDelay = duration
			}
		}
	}
	if totalReceiveCount == 0 {
		return
	}
	avgDelay = totalDelay / time.Duration(totalReceiveCount)
	return
}
