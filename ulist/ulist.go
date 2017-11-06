package ulist

import (
	"os"

	"github.com/inconshreveable/log15"
	"github.com/naoina/toml"
	"github.com/pkg/errors"
)

const tomlListPath = ".cshelper/ulist.toml"

var (
	allUsers = Ulist{}
	Log      = log15.New()
)

type Ulist struct {
	users []User
}

func NewUlist(users []User) *Ulist {
	return &Ulist{users: users}
}

type tomlUlist struct {
	Users []User `toml:"user,omitempty"`
}

func Init(csvListPath string, ignoreUlistToml bool) error {
	// parse ulist from google form export
	users, err := ParseFile(csvListPath)
	if err != nil {
		return errors.Wrapf(err, "unable to parse %q", csvListPath)
	}
	allUsers = *NewUlist(users)

	if ignoreUlistToml {
		return err
	}
	if _, err := os.Stat(tomlListPath); os.IsNotExist(err) {
		return allUsers.Save(tomlListPath)
	}
	return allUsers.Load(tomlListPath)
}

func GetAllUsers() *Ulist {
	return &allUsers
}

func (ul *Ulist) Len() int {
	return len(ul.users)
}

func (ul *Ulist) Iter(iterFn func(i int, u User) error) error {
	for i, u := range ul.users {
		if err := iterFn(i, u); err != nil {
			return err
		}
	}
	return nil
}

func (ul *Ulist) Get(i int) *User {
	return &ul.users[i]
}

func (ul *Ulist) GetByName(name string) *User {
	for _, u := range ul.users {
		if CompareUsernames(u.Username, name) {
			return &u
		}
	}
	return nil
}

func (ul *Ulist) Filter(filterFn func(u User) bool) *Ulist {
	var new []User
	ul.Iter(func(_ int, u User) error {
		if filterFn(u) {
			new = append(new, u)
		}
		return nil
	})
	return NewUlist(new)
}

func (ul *Ulist) Save(fpath string) error {
	Log.Debug("saving toml", "path", fpath)
	// save ulist as toml
	fd, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "unable to open %q for writing", fpath)
	}
	enc := toml.NewEncoder(fd)
	if err := enc.Encode(&tomlUlist{Users: ul.users}); err != nil {
		return errors.Wrap(err, "unable to encode user list as toml")
	}
	return nil
}

func (ul *Ulist) Load(fpath string) error {
	Log.Debug("loading toml", "path", fpath)
	// decode ulist from ulist.toml
	fd, err := os.Open(fpath)
	if err != nil {
		return errors.Wrapf(err, "unable to open %q for reading", fpath)
	}
	dec := toml.NewDecoder(fd)
	var tomlData tomlUlist
	if err := dec.Decode(&tomlData); err != nil {
		return errors.Wrap(err, "unable to decode user list toml")
	}
	*ul = Ulist{users: tomlData.Users}
	return nil
}
