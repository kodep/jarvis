package factories

import (
	"github.com/kodep/jarvis/internal/mattermost/handlers"
)

func SendTyping(client handlers.WSClient, channelID string) {
	const action = "user_typing"

	client.SendMessage(action, map[string]interface{}{
		"channel_id": channelID,
	})
}
