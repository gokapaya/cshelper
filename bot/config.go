package bot

type Config struct {
	Useragent    string
	ClientID     string
	ClientSecret string
	Username     string
	Password     string
}

const DefaultSubject = "/r/ClosetSanta notification"
const footnote = `
---
This bot is maintained by /u/gokapaya  
The source code is available at [github](https://github.com/gokapaya/cshelper)
`
