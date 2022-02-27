package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	helloworldpb "server/proto/helloworld"

	"github.com/go-redis/redis"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func handleWebSocket(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		panic(err)
	}
	go func() {
		defer conn.Close()

		for {
			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				errClose := wsutil.ClosedError{Code: 1000}
				if strings.EqualFold(err.Error(), errClose.Error()) {
					println("ws client close")
					return
				}
				panic(err)
			}

			println("msg: ", msg)
			println("op: ", op)

			event := helloworldpb.TaskEvent{}
			err = json.Unmarshal(msg, &event)
			if err != nil {
				fmt.Println(err)
				continue
			}
			go handleTaskEvent(conn, op, &event)
		}
	}()
}

func handleTaskEvent(conn net.Conn, op ws.OpCode, event *helloworldpb.TaskEvent) {
	fmt.Println(event)
	switch event.GetOp() {
	case helloworldpb.TaskEvent_STATUS.String():
		taskIds := event.GetTaskIds()
		if len(taskIds) == 0 {
			return
		}
		mem := make([]string, len(taskIds))
		for {
			for idx, taskId := range taskIds {
				k := fmt.Sprintf("task-status-%d", taskId)
				v, err := rdb.Get(context.TODO(), k).Result()
				if err == redis.Nil {
					continue
				} else if err != nil {
					fmt.Println(err)
				} else {
					mem[idx] = v
					resp := helloworldpb.TaskEventResponse{
						TaskId: taskId,
						Status: helloworldpb.
							TaskEventResponse_TaskStatus(
								helloworldpb.TaskEventResponse_TaskStatus_value[v]),
					}
					msg, err := json.Marshal(&resp)
					if err != nil {
						fmt.Println(err)
						continue
					}
					err = wsutil.WriteServerMessage(conn, op, msg)
					if err != nil {
						errClose := wsutil.ClosedError{Code: 1000}
						if strings.EqualFold(err.Error(), errClose.Error()) {
							fmt.Println("ws client close")
							return
						}
						fmt.Println(err)
					}
				}
			}
			time.Sleep(1 * time.Second)
		}
	default:
		fmt.Println("unknown task op")
	}
}
