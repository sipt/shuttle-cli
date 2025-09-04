package apis

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	red    = color.RGB(255, 66, 69)
	green  = color.RGB(48, 209, 88)
	white  = color.RGB(255, 255, 255)
	orange = color.RGB(255, 146, 48)
)

func PrintRtt(rtt string) {
	switch rtt {
	case "no rtt", "failed":
		red.Print(rtt)
		return
	case "":
		white.Print("unknown")
		return
	default:
	}

	// rtt is 321ms
	// color:
	// green : rtt<300ms
	// yellow : 300ms<=rtt<600ms
	// red : rtt>=600ms
	rttValue := 0
	fmt.Sscanf(rtt, "%d", &rttValue)

	if rttValue < 300 {
		green.Print(rtt)
	} else if rttValue < 600 {
		orange.Print(rtt)
	} else {
		red.Print(rtt)
	}
}
