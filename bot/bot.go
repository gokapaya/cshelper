package bot

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gokapaya/cshelper/ulist"
	"github.com/inconshreveable/log15"
	"github.com/pkg/errors"
	"github.com/turnage/graw/reddit"
)

var (
	bot CSBot
	Log log15.Logger
)

type CSBot struct {
	reddit.Bot
}

func Init(cfg *Config) error {
	botHandle, err := reddit.NewBot(reddit.BotConfig{
		Agent: cfg.Useragent,
		App: reddit.App{
			ID:       cfg.ClientID,
			Secret:   cfg.ClientSecret,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})
	if err != nil {
		return errors.Wrap(err, "creating bot handle failed")
	}

	bot = CSBot{
		Bot: botHandle,
	}
	Log.Debug("bot instance created", "username", cfg.Username)
	return nil
}

func PmUserWithTemplate(user ulist.User, subject string, t *template.Template) error {
	var (
		body string
	)

	var buf bytes.Buffer
	if err := t.Execute(&buf, user); err != nil {
		Log.Debug("dumping template", "t", t)
		Log.Debug("dumping user", "u", user)
		return errors.Wrap(err, "rendering template failed")
	}

	if _, err := buf.WriteString(footnote); err != nil {
		return errors.Wrap(err, "writing footnote failed")
	}
	body = buf.String()

	Log.Debug("sending message", "re", subject)
	return bot.SendMessage(user.Username, subject, body)
}

type Subreddit struct {
	rawName string
}

func NewSubreddit(name string) Subreddit {
	return Subreddit{rawName: name}
}
func (sub Subreddit) Name() string {
	return fmt.Sprintf("/r/%s", sub.rawName)
}
func PmSubredditWithTemplate(sub Subreddit, subject string, t *template.Template) error {
	var (
		body string
		err  error
	)

	body, err = renderTemplate(t, sub)
	if err != nil {
		return err
	}

	Log.Debug("sending message", "re", subject)
	return bot.SendMessage(sub.Name(), subject, body)
}

func renderTemplate(t *template.Template, data interface{}) (string, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		Log.Debug("dumping template", "t", t)
		Log.Debug("dumping data", "data", data)
		return "", errors.Wrap(err, "rendering template failed")
	}

	if _, err := buf.WriteString(footnote); err != nil {
		return "", errors.Wrap(err, "writing footnote failed")
	}
	return buf.String(), nil
}
