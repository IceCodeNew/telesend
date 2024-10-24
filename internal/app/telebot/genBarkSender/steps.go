package genBarkSender

import (
	"fmt"
	"strings"
)

var (
	nextStep = 1

	steps = []string{
		"set the DeviceKey",
		"set the Bark Server address",
		"request a string for verifying the Bark Sender",
		"finish the setup",
	}
)

func previousStepsNotComplete(currStep int) string {
	if currStep < 2 {
		return ""
	}

	var _reply strings.Builder
	_reply.WriteString(
		fmt.Sprintf("ERROR: [Internal] Steps before %d did not complete successfully:", currStep),
	)
	for i := 0; i < currStep-1; i++ {
		_reply.WriteString(
			fmt.Sprintf("\n   Step %d: %s", i, steps[i]),
		)
	}

	reply := _reply.String()
	return reply
}
