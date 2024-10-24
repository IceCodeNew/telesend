package telebot

import (
	"fmt"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func payloadsOfMessage(ctx *ext.Context) []string {
	// Syntax: "</command>@<bot> <payload>"
	_input := strings.Fields(ctx.EffectiveMessage.Text)
	if len(_input) < 2 {
		return nil
	}
	return _input[1:]
}

func replyNoDetailInternalErr(bot *gotgbot.Bot, ctx *ext.Context, upstreamErr error, reply string) error {
	_, _ = ctx.EffectiveMessage.Reply(bot, reply, nil)

	if upstreamErr != nil {
		return fmt.Errorf(`
%s
DEBUG: original error was:
%v`,
			reply, upstreamErr,
		)
	}
	return fmt.Errorf(reply)
}

func previousStepsNotComplete(currStep int, steps []string) string {
	if currStep < 2 {
		return ""
	}

	var _reply strings.Builder
	_reply.WriteString(fmt.Sprintf(
		"ERROR: [Internal] Steps before %d did not complete successfully:", currStep,
	))
	for i := 0; i < currStep-1; i++ {
		_reply.WriteString(fmt.Sprintf(
			"\n   Step %d: %s", i, steps[i],
		))
	}

	reply := _reply.String()
	return reply
}
