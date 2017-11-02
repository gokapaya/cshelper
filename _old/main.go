package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("closetsanta", "help /r/ClosetSanta")
	dbpath := app.StringOpt("db", "./giftee.db", "path to the giftee database")
	useragent := app.StringOpt("useragent", "./useragent.protobuf", "path to the reddit bot useragent")
	templatedir := app.StringOpt("templates", "./templates", "path to the templates directory")

	a, err := InitApp(*useragent, *templatedir, *dbpath)
	if err != nil {
		log.Fatal(err)
	}

	if err := a.db.Open(); err != nil {
		log.Fatal(err)
	}
	defer a.db.Close()

	app.Command("parse-csv", "get the users from the Google docs CSV", func(c *cli.Cmd) {
		dry := c.BoolOpt("dry", false, "just print contents of csv")

		csvpath := c.StringArg("CSV", "", "path to the CSV file")

		c.Action = func() {
			if err := a.ParseCSV(*dry, *csvpath); err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("parse-shipping-csv", "get shipping status from csv", func(c *cli.Cmd) {
		dry := c.BoolOpt("dry", false, "just print contents of csv")

		csvpath := c.StringArg("CSV", "", "path to the CSV file")

		c.Action = func() {
			if err := a.ParseCSVwithShippingStatus(*dry, *csvpath); err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("list-users", "list all user in the database", func(c *cli.Cmd) {
		c.Action = func() {
			err := a.ListUsers()
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("print-user", "print the data about a single user", func(c *cli.Cmd) {
		user := c.StringArg("USER", "", "user to print out")

		c.Action = func() {
			err := a.PrintUser(*user)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("cleanup", "clean up data in the database (mainly for country names)", func(c *cli.Cmd) {
		c.Action = func() {
			err := a.db.CleanDatabase()
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("match", "match users to each other", func(c *cli.Cmd) {
		c.Action = func() {
			if err := a.GenerateMatchList(); err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("match-export", "Save a CSV file with the matched pairs", func(c *cli.Cmd) {
		dry := c.BoolOpt("dry", false, "don't run live")
		filename := c.StringArg("FILENAME", "", "the filename to export to")

		c.Action = func() {
			if err := a.ExportMatchList(*dry, *filename); err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("pm-user", "send a PM to a user", func(c *cli.Cmd) {
		c.LongDesc = "Using this command without the template option you can enter subject and content yourself without defining a template."

		dry := c.BoolOpt("dry", false, "don't run live")
		withTemplate := c.StringOpt("with-template", "", "send a message generated from a template")

		user := c.StringArg("USER", "", "username to contact")
		subject := c.StringArg("SUBJECT", "", "the subject of the PM")

		c.Before = func() {
			if *user == "" {
				fmt.Println("Needs a username")
				os.Exit(3)
			}

			if *subject == "" {
				fmt.Println("Needs a subject")
				os.Exit(3)
			}
		}

		c.Action = func() {
			err := a.RedditMessageTo(*dry, *subject, *withTemplate, *user)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	// app.Command("pm-user-csv", "send a PM to a user based on CSV information", func(c *cli.Cmd) {

	// 	dry := c.BoolOpt("dry", false, "don't run live")
	// 	withTemplate := c.StringOpt("with-template", "", "send a message generated from a template")
	// 	withCSVPath := c.StringOpt("with-csv", "matches.csv", "path to the csv file")

	// 	user := c.StringArg("USER", "", "username to contact")
	// 	subject := c.StringArg("SUBJECT", "", "the subject of the PM")

	// 	c.Before = func() {
	// 		if *user == "" {
	// 			fmt.Println("Needs a username")
	// 			os.Exit(3)
	// 		}

	// 		if *subject == "" {
	// 			fmt.Println("Needs a subject")
	// 			os.Exit(3)
	// 		}

	// 		if *withTemplate == "" {
	// 			fmt.Println("Needs a template (for now)")
	// 			os.Exit(3)
	// 		}
	// 	}

	// 	c.Action = func() {
	// 		err := a.RedditMessageFromCSV(*dry, *withCSVPath, *subject, *withTemplate, *user)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	// })

	app.Command("pm-batch", "send PMs to users", func(c *cli.Cmd) {
		c.LongDesc = "Using this command without the template option you can enter subject and content yourself without defining a template."

		dry := c.BoolOpt("dry", false, "don't run live")
		withTemplate := c.StringOpt("with-template", "", "send a message generated from a template")

		subject := c.StringArg("SUBJECT", "", "the subject of the PM")

		c.Before = func() {
			if *subject == "" {
				fmt.Println("Needs a subject")
				os.Exit(3)
			}
		}

		c.Action = func() {
			err := a.RedditMessageToAll(*dry, *subject, *withTemplate)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("pm-batch-csv", "send PMs to users from a Pair CSV", func(c *cli.Cmd) {
		c.LongDesc = "Using this command without the template option you can enter subject and content yourself without defining a template."

		dry := c.BoolOpt("dry", false, "don't run live")

		csvpath := c.StringArg("PATH", "", "the path to the csv")

		c.Action = func() {
			var withTemplate = ""
			var subject = ""
			err := a.RedditMessageFromCSV(*dry, *csvpath, subject, withTemplate)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("pm-batch-rematch", "send rematch PMs to users", func(c *cli.Cmd) {
		dry := c.BoolOpt("dry", false, "don't run live")

		c.Action = func() {
			err := a.MessageUsersWithRematchPM(*dry)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Command("pm-batch-shipping", "send shipping status PMs to users", func(c *cli.Cmd) {
		dry := c.BoolOpt("dry", false, "don't run live")

		csvpath := c.StringArg("CSV", "", "path to the CSV file")

		c.Action = func() {
			err := a.MessageUsersWithShippingStatusPM(*dry, *csvpath)
			if err != nil {
				log.Fatal(err)
			}
		}
	})

	app.Run(os.Args)
}
