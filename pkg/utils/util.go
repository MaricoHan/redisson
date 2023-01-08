package utils

import (
	"runtime"
	"strconv"
	"strings"
)

// GoID
// @Description: get the id of current goroutine
// @return int64
func GoID() int64 {
	buf := make([]byte, 35)
	runtime.Stack(buf, false)
	s := string(buf)
	parseInt, _ := strconv.ParseInt(strings.TrimSpace(s[10:strings.IndexByte(s, '[')]), 10, 64)
	return parseInt
}

func ChannelName(name string) string {
	return "redisson_lock__channel" + ":{" + name + "}"
}
