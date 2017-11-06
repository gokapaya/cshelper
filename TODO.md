# TODO

List of stuff that needs to be ported over from the old implementation.

- [x] parse google doc export on init
- [x] keep original export as backup
- [x] write toml list for emergency edits

- [x] parse msg templates on init
  - [ ] custom frontmatter for reddit 'regarding' field
- [ ] review templates:
  - [ ] welcome to the contest/match msg
    - ship as gift
    - explain message bot
  - [ ] gift on the way msg
  - [ ] sorry no ship msg
  - [ ] rematcher welcome

- [ ] matching
  - [x] look at munkres algorithm http://csclab.murraystate.edu/~bob.pilgrim/445/munkres.html
  - [x] improve overseas / non-overseas algorithms
  - [x] create/export match list for RSSB (2-2 csv)
  - [ ] do not match with user from same represented sub

- [ ] reddit bot
  - [x] setup bot
  - [x] pm single user with msg template
  - [ ] pm batch with template

- [x] additional documents
  - [x] template: asking subs for participation
  - [x] list of subs to contact

### subcommands:

- pm
  - needs to load ulist/templates
  - needs to setup bot
  - subcommands:
    - single USER -> sends pm to user, choose template with flag (-t)
    - batch [USER,...] -> sends batch pms to all or listed
    - ["all", "rematch"] -> special words

- match
  - needs to load ulist
  - matching...
  - output 2-2 csv
  - annotate in toml
