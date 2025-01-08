package consensus

import (
	"time"
)

func LeaderNode(firstBlockTime int64, step int64, totalNodes int) int {
	currentTime := time.Now().Unix()
	return int(((currentTime - firstBlockTime) / step) % int64(totalNodes))
}
