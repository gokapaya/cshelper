# cshelper

### WARNING! Under construction.

`cshelper` helps us organize the secret santa over on [/r/ClosetSanta](https://reddit.com/r/closetsanta).

## Features

```
Usage: closetsanta [OPTIONS] COMMAND [arg...]

Options:
  --db="./giftee.db"                   path to the giftee database
  --useragent="./useragent.protobuf"   path to the reddit bot useragent
  --templates="./templates"            path to the templates directory

Commands:
  parse-csv            get the users from the Google docs CSV
  parse-shipping-csv   get shipping status from csv
  list-users           list all user in the database
  print-user           print the data about a single user
  cleanup              clean up data in the database (mainly for country names)
  match                match users to each other
  match-export         Save a CSV file with the matched pairs
  pm-user              send a PM to a user
  pm-batch             send PMs to users
  pm-batch-csv         send PMs to users from a Pair CSV
  pm-batch-rematch     send rematch PMs to users
  pm-batch-shipping    send shipping status PMs to users
```
