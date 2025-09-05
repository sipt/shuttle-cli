package apis

import (
	"time"

	"github.com/fatih/color"
)

var (
	red    = color.RGB(255, 66, 69)
	green  = color.RGB(48, 209, 88)
	white  = color.RGB(255, 255, 255)
	orange = color.RGB(255, 146, 48)
)

func PrintRtt(rtt int64) {
	switch {
	case rtt == 0:
		red.Print("no rtt")
		return
	case rtt < 0:
		orange.Print("failed")
		return
	default:
	}

	// rtt is 321ms
	// color:
	// green : rtt<300ms
	// yellow : 300ms<=rtt<600ms
	// red : rtt>=600ms

	rttValue := time.Duration(rtt) * time.Millisecond
	if rtt < 300 {
		green.Print(rttValue.String())
	} else if rtt < 600 {
		orange.Print(rttValue.String())
	} else {
		red.Print(rttValue.String())
	}
}
