package timer

import (
	"go/cmkj_server_go/util"

	"time"
)

var GameTimer map[int]int

func StartGameTimer() {
	GameTimer = make(map[int]int)
	trick := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-trick.C:
			//util.Log.Debug("game timer")
			for k, t := range GameTimer {
				util.Log.Debug(k, t)
				if t > 0 {
					GameTimer[k] = t - 1
				}
			}
		}
	}
}

func init() {
	go StartGameTimer()
}
