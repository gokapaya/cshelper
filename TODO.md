# TODO

List of stuff that needs to be ported over from the old implementation.

- [ ] parse google doc export on init
- [ ] keep original export as backup
- [ ] write toml list for emergency edits

- [ ] parse msg templates on init
  - [ ] custom frontmatter for reddit 'regarding' field
- [ ] review templates:
  - [ ] welcome to the contest
  - [ ] matched msg
  - [ ] shipping status msg
  - [ ] sorry no ship msg
  - [ ] rematcher welcome

- [ ] new templates
  - [ ] asking subs for participation

- [ ] matching
  - [ ] look at munkres algorithm http://csclab.murraystate.edu/~bob.pilgrim/445/munkres.html
  - [ ] improve overseas / non-overseas algorithms
  - [ ] create/export match list for ClosetSantaMessageBot (2-2 csv)
  - [ ] do not match with user from same represented sub

- [-] reddit bot
  - [ ] setup bot
  - [ ] pm single user with msg template
  - [ ] pm batch with template

** subcommands:

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
