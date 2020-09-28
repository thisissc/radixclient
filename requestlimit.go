package radixclient

import (
	"fmt"
	"log"
	"time"

	"github.com/mediocregopher/radix/v3"
)

//Usage:
//if radixclient.RequestLimit("CompleteMyTaskHandler", userid, 1, 1) {
//	return h.RenderMsg(400, 400, "Too Fast")
//}
func RequestLimit(zone, member string, duration, burst uint) bool {
	if duration < 1 {
		duration = 1
	}

	c := Radix()

	timeFlag := time.Now().Unix() / int64(duration)
	mutexKey := fmt.Sprintf("ZA:REQUESTLIMIT:%s:%s:%d", zone, member, timeFlag)
	var cnt int
	err := c.Do(radix.Cmd(&cnt, "INCR", mutexKey))
	if err != nil {
		log.Println(err)
	}
	if cnt == 1 {
		c.Do(radix.FlatCmd(nil, "EXPIRE", mutexKey, duration+60))
	}

	return cnt > int(burst)
}
