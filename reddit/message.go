package reddit

import (
	"bytes"
	"errors"
	"html/template"
	"log"

	"github.com/gokapaya/cshelper/csv"
	"github.com/gokapaya/cshelper/errlist"
	"github.com/gokapaya/cshelper/match"
)

// SendCustomMessageTo ...
func (s *SantaBot) SendCustomMessageTo(subject, content string, ul ...csv.User) error {
	var el errlist.ErrList

	if subject == "" {
		return errors.New("missing subject")
	}
	if content == "" {
		return errors.New("missing content")
	}

	for _, u := range ul {
		log.Printf("messaging... /u/%v :: %v\n", u.Username, subject)
		if err := s.bot.SendMessage(u.Username, subject, content); err != nil {
			el = append(el, err)
		}
	}

	if el.NotEmpty() {
		return el
	}
	return nil
}

// SendCustomTemplateMessageTo ...
func (s *SantaBot) SendCustomTemplateMessageTo(subject string, t template.Template, data interface{}, ul ...csv.User) error {
	var el errlist.ErrList

	if subject == "" {
		return errors.New("missing subject")
	}

	for _, u := range ul {
		var content = new(bytes.Buffer)
		if err := t.Execute(content, data); err != nil {
			el = append(el, err)
		}

		log.Printf("messaging... /u/%v :: %v\n", u.Username, subject)
		if err := s.bot.SendMessage(u.Username, subject, content.String()); err != nil {
			el = append(el, err)
		}
	}

	if el.NotEmpty() {
		return el
	}
	return nil
}

// SendTemplateMessageTo ...
func (s *SantaBot) SendTemplateMessageTo(subject string, t template.Template, pl ...match.Pair) error {
	var el errlist.ErrList

	if subject == "" {
		return errors.New("missing subject")
	}

	for _, p := range pl {
		var content = new(bytes.Buffer)
		if err := t.Execute(content, p.Giftee); err != nil {
			el = append(el, err)
		}

		log.Printf("messaging... /u/%v :: %v\n", p.Santa.Username, subject)
		if err := s.bot.SendMessage(p.Santa.Username, subject, content.String()); err != nil {
			el = append(el, err)
		}
	}

	if el.NotEmpty() {
		return el
	}
	return nil
}
