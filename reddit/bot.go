package reddit

import (
	"time"

	"github.com/turnage/graw/reddit"
)

// SantaBot wraps reddit.Bot
type SantaBot struct {
	bot reddit.Bot
}

// NewSantaBot returns a new SantaBot with ours config
func NewSantaBot(agent string) (*SantaBot, error) {
	b, err := reddit.NewBotFromAgentFile(agent, 5*time.Second)
	if err != nil {
		return nil, err
	}

	return &SantaBot{bot: b}, nil
}
