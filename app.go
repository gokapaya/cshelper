package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/gokapaya/cshelper/boltclient"
	"github.com/gokapaya/cshelper/csv"
	"github.com/gokapaya/cshelper/match"
	"github.com/gokapaya/cshelper/reddit"
	"github.com/gokapaya/cshelper/templates"
)

// App pulls together submodule functionality
type App struct {
	// bot is an instance of the reddit bot lib
	bot reddit.SantaBot
	// db is an instance of the database client
	db boltclient.Client
	// templates is a list of message templates
	templates []template.Template
}

// InitApp opens the database client and the bot client,
// parses the template directory and returns an App instance
func InitApp(useragent, templateDir, dbPath string) (*App, error) {
	b, err := reddit.NewSantaBot(useragent)
	if err != nil {
		return nil, err
	}
	t, err := templates.ParseDir(templateDir)
	if err != nil {
		return nil, err
	}
	log.Println(len(t), "template(s) found")
	d := boltclient.NewClient(dbPath)
	if err != nil {
		return nil, err
	}

	return &App{
		bot:       *b,
		templates: t,
		db:        *d,
	}, err
}

// ParseCSV takes the path to a CSV from the Google Form.
// It saves the user data in the database
func (a *App) ParseCSV(dry bool, path string) error {
	userList, err := csv.GetUserList(path)
	if err != nil {
		return err
	}
	if dry {
		for i, u := range userList {
			fmt.Printf("%3v - %v\n", i, u.Username)
		}
	} else {
		return a.db.StoreUsers(userList...)
	}

	return nil
}

// ParseCSVwithShippingStatus ...
func (a *App) ParseCSVwithShippingStatus(dry bool, path string) error {
	stmap, err := csv.GetShippingStatusInf(path)
	if err != nil {
		return err
	}

	log.Println(len(stmap))

	// dirty code to deal with dirty DB
	var nopresent = func(status string, pair *match.Pair) bool {
		switch pair.Giftee.Username {
		case "gokapaya":
			return false
		case "Curbadon":
			return false
		case "Jaffasaurs":
			return false
		case "ShinobuOwO":
			return false
		case "Victinithetiny101":
			return false
		default:
			if status == "Failure" {
				return true
			}
			return false
		}
	}

	for u, s := range stmap {
		if s == "Resolved Failure" {
			continue
		}
		p, err := a.db.GetPair(u)
		if err != nil {
			log.Println(u, s)
			return err
		}
		// filter out the failures
		if nopresent(s, p) {
			// log.Printf("%20v -x-> %v\n", p.Santa.Username, p.Giftee.Username)
			fmt.Printf("%v\n", p.Giftee.Username)
		}
		// filter out the to be checked
		if s == "Check" {
			// log.Printf("%20v -!-> %v\n", p.Santa.Username, p.Giftee.Username)
		}
	}

	return nil
}

// MessageUsersWithRematchPM sends the rematch message to all rematchers
func (a *App) MessageUsersWithRematchPM(dry bool) error {
	ul, err := a.db.GetUserList()
	if err != nil {
		return err
	}
	// load the content of the rematch message
	content, err := ioutil.ReadFile("./templates/rematch.templ")
	if err != nil {
		return err
	}

	var regifters []csv.User
	for _, u := range ul {
		if u.Regift {
			regifters = append(regifters, u)
		}
	}

	log.Println(len(regifters))

	if !dry {
		if err := a.bot.SendCustomMessageTo("/r/ClosetSanta - Rematch", string(content), regifters...); err != nil {
			return err
		}
	} else {
		for _, u := range regifters {
			log.Println("messaging...", u.Username, u.Regift)
		}
	}

	return nil
}

// MessageUsersWithShippingStatusPM sends their gifts status to all users
// needs a csv from google docs with format:
// santa, status
func (a *App) MessageUsersWithShippingStatusPM(dry bool, path string) error {
	stmap, err := csv.GetShippingStatusInf(path)
	if err != nil {
		return err
	}
	log.Println(len(stmap))

	// load the shipping status template
	t, err := templates.Find("shipping.templ", a.templates)
	if err != nil {
		return err
	}

	// dirty code to deal with dirty DB
	var nopresent = func(status string, pair *match.Pair) bool {
		switch pair.Giftee.Username {
		case "gokapaya":
			return false
		case "Curbadon":
			return false
		case "Jaffasaurs":
			return false
		case "ShinobuOwO":
			return false
		case "Victinithetiny101":
			return false
		default:
			if status == "Failure" {
				return true
			}
			return false
		}
	}
	var getspresent = func(status string, pair *match.Pair) bool {
		switch pair.Giftee.Username {
		case "gokapaya":
			return false
		case "Curbadon":
			return false
		case "Jaffasaurs":
			return false
		case "ShinobuOwO":
			return false
		case "Victinithetiny101":
			return false
		default:
			if status == "Yes" {
				return true
			}
			return false
		}
	}

	var failure, success []csv.User
	for u, status := range stmap {
		if status == "Resolved Failure" || status == "Check" {
			continue
		}
		p, err := a.db.GetPair(u)
		if err != nil {
			log.Println(u, status)
			return err
		}
		if nopresent(status, p) {
			failure = append(failure, p.Giftee)
		}
		if getspresent(status, p) {
			success = append(success, p.Giftee)
		}
	}

	if !dry {
		// fmt.Println("-- SUCCESS --")
		// if err := a.bot.SendCustomTemplateMessageTo("/r/ClosetSanta - Shipping Status", *t, true, success...); err != nil {
		// 	return err
		// }
		fmt.Println("-- FAILURE --")
		if err := a.bot.SendCustomTemplateMessageTo("/r/ClosetSanta - Shipping Status", *t, false, failure...); err != nil {
			return err
		}
	} else {
		log.Println("successes")
		for _, u := range success {
			log.Println("messaging...", u.Username)
		}
		t.Execute(os.Stdout, true)
		log.Println("failures")
		for _, u := range failure {
			log.Println("messaging...", u.Username)
		}
		t.Execute(os.Stdout, false)
	}

	log.Println("failures:", len(failure))
	log.Println("successes:", len(success))

	return nil
}

// ListUsers prints all usernames in the database
func (a *App) ListUsers() error {
	userList, err := a.db.GetUserList()
	if err != nil {
		return err
	}
	var regifters int
	var germans []string
	for i, u := range userList {
		if u.Regift {
			regifters++
		}
		fmt.Printf("%3v - %20v - %20v - %v\n", i, u.Username, u.RepSubreddit, u.Address.Country)
	}

	fmt.Printf("we have %v potential rematchers\n", regifters)
	fmt.Println(germans)

	return nil
}

// PrintUser prints all the data on user `name` from the database
func (a *App) PrintUser(name string) error {
	u, err := a.db.GetUser(name)
	if err != nil {
		return err
	}
	fmt.Printf("/u/%v ~ %v\nshare username? %v\nships international? %v\nrematcher? %v\n--- address:\n%v\n--- message to santa:\n%v\n--- list:\n%v\n",
		u.Username, u.RepSubreddit, u.ShareName, u.International, u.Regift, u.Address.String(), u.MessageToSanta, u.Watchlist)

	return nil
}

// RedditMessageTo takes a name and an optional templateName
// Sends a PM to user `name` from database
// If templateName is given it uses the template
// otherwise a message has to be provided
func (a *App) RedditMessageTo(dry bool, subject, templateName, name string) error {
	// user, err := a.db.GetUser(name)
	pair, err := a.db.GetPair(name)
	if err != nil {
		return err
	}
	// without template
	if templateName == "" {
		// ask for subject and content
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter message: ")
		content, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Println()

		if dry {
			fmt.Fprint(os.Stdout, "running dry...")
			fmt.Fprintf(os.Stdout, "contacting %v\nwith subject: <%v>\n---\n%v---\n", pair.Santa.Username, subject, content)
			return nil
		}
		return a.bot.SendCustomMessageTo(subject, content, pair.Santa)
	}

	// with template
	fmt.Fprintf(os.Stdout, "searching for '%v' in templates...\n", templateName)
	template, err := templates.Find(templateName, a.templates)
	if err != nil {
		return err
	}

	if dry {
		fmt.Fprintln(os.Stdout, "running dry...")
		fmt.Fprintf(os.Stdout, "contacting %v\nwith subject: %v\nusing template: %v\n", pair.Santa.Username, subject, template.Name())
		return template.Execute(os.Stdout, pair.Giftee)
	}

	return a.bot.SendTemplateMessageTo(subject, *template, *pair)
}

// RedditMessageToAll is the same as above
// for all users in database
func (a *App) RedditMessageToAll(dry bool, subject, templateName string) error {
	// userList, err := a.db.GetUserList()
	pairList, err := a.db.GetMatches()
	if err != nil {
		return err
	}
	// without template
	if templateName == "" {
		// ask for subject and content
		// reader := bufio.NewReader(os.Stdin)
		// fmt.Print("Enter message: ")
		// content, err := reader.ReadString('\n')
		// if err != nil {
		// 	return err
		// }

		if dry {
			for _, pair := range pairList {
				fmt.Fprintln(os.Stdout, "running dry...")
				fmt.Fprintf(os.Stdout, "contacting %v\nwith subject: %v\n", pair.Santa.Username, subject)
			}
			return nil
		}
		log.Fatal("not implemented")
		// return a.bot.SendCustomMessageTo(subject, content, pairList...)
	}
	// with template
	template, err := templates.Find(templateName, a.templates)
	if err != nil {
		return err
	}

	if dry {
		for _, pair := range pairList {
			fmt.Fprintln(os.Stdout, "running dry...")
			fmt.Fprintf(os.Stdout, "contacting %v\nwith subject: %v\nusing template: %v\n", pair.Santa.Username, subject, template.Name())
		}
		return nil
	}

	return a.bot.SendTemplateMessageTo(subject, *template, pairList...)
}

// RedditMessageFromCSV sends messages to people on a SantaList.csv
// Used for reaching out to people during the contest
func (a *App) RedditMessageFromCSV(dry bool, csvPath, subject, templateName string) error {
	pairs, err := csv.GetPairs(csvPath)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, "running dry...")

	var messages = 0
	for santaName := range pairs {
		log.Println(">>> now:", santaName, "==>", pairs[santaName])
		giftee, err := a.db.GetUser(pairs[santaName])
		if err != nil {
			return err
		}
		santa, err := a.db.GetUser(santaName)
		if err != nil {
			return err
		}

		pairRegifter := match.Pair{
			Santa:  *santa,
			Giftee: *giftee,
		}
		pairMessageGiftee := match.Pair{
			Santa:  *giftee,
			Giftee: *santa,
		}

		templateRegifter, err := templates.Find("regifter-msg.templ", a.templates)
		if err != nil {
			return err
		}
		templateGiftee, err := templates.Find("giftee-regift-msg.templ", a.templates)
		if err != nil {
			return err
		}

		var subject1 = "/u/ClosetSanta - You are a regifter!"
		var subject2 = "/u/ClosetSanta - You have been rematched!"

		if dry {
			fmt.Fprintf(os.Stdout, "contacting \t %10v\nwith subject: \t %10v\nusing template: \t %10v\n", pairRegifter.Santa.Username, subject1, templateRegifter.Name())
			messages++
			fmt.Fprintf(os.Stdout, "contacting \t %10v\nwith subject: \t %10v\nusing template: \t %10v\n", pairMessageGiftee.Santa.Username, subject2, templateGiftee.Name())
			messages++
			// template.Execute(os.Stdout, pair.Giftee)
			// return nil
			continue
		}

		if err := a.bot.SendTemplateMessageTo(subject1, *templateRegifter, pairRegifter); err != nil {
			log.Print(err)
		}
		if err := a.bot.SendTemplateMessageTo(subject2, *templateGiftee, pairMessageGiftee); err != nil {
			log.Print(err)
		}
		log.Println("-------------")
	}

	log.Println(messages, "/", 30)

	return nil
}

// GenerateMatchList uses the database and matches users in pairs
// It generates a CSV file with all matches
// and reports matching information
func (a *App) GenerateMatchList() error {
	userList, err := a.db.GetUserList()
	if err != nil {
		return err
	}
	matchList, err := match.MatchUsers(userList)
	if err != nil {
		return err
	}
	// write matchList to database
	if err := a.db.StoreMatches(matchList...); err != nil {
		return err
	}
	// write matchList to CSV
	if err := match.WriteToCSV(true, matchList, ""); err != nil {
		return err
	}
	// for _, p := range matchList {
	// 	fmt.Printf("%v --> %v\n", p.Santa.Username, p.Giftee.Username)
	// }
	fmt.Printf("matched %v of %v pairs\n", len(matchList), len(userList))
	return nil
}

// ExportMatchList writes the matched pairs to a CSV file with the SantaList.csv formating
// used by the anonymous message bot
func (a *App) ExportMatchList(dry bool, filename string) error {
	matchList, err := a.db.GetMatches()
	if err != nil {
		return err
	}
	// write matchList to CSV
	if err := match.WriteToCSV(dry, matchList, filename); err != nil {
		return err
	}
	if dry {
		fmt.Println("dry run... Not writing to file")
		return nil
	}
	fmt.Printf("Exported match list to %v\n", filename)
	return nil
}
