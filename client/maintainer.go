package client

import (
	"time"
)

func (cl *Client) maintainSession() {
	cl.Log.Println("Maintaining session")
	for cl.isRunning.Load() == 1 {
		cl.Log.Println("Next Heartbeat: ", time.Now().Add(cl.HeartbeatInterval).Format(time.DateTime))
		select {
		case <-time.After(cl.HeartbeatInterval):
			if err := cl.sendHeartbeat(); err != nil {
				cl.Log.Printf("Heartbeat failed: %v", err)
				cl.isRunning.Store(0)
				return
			}
			cl.Log.Println("Heartbeat Sent")
		case <-cl.Ctx.Done():
			return
		}
	}
}
