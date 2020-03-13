# BST WEB

WIP web-server leveraging Auth0 or another identity provider.

Build for *nix with:

```
env GOOS=linux GOARCH=amd64 go build
``` 

Or for your own OS with:
```
go build
```

Run with:

```
./bst_web \
    -static="./dist" \
    -index="index.html" \
    -404="404.html" \
    -js="/js" \
    -media="/media" \
    -public="/public" \
    -protected="/protected" \
    -clientid="clientid" \
    -clientsecret="clientsecret" \
    -issuer="https://issuer.com/" \
    -audience="myaudience" \
    -callback="/callback" \
    -filestorekey="averysecurekey" \
    -host="my.host.com" \
    -port="443"
```

---

## To-do

### Auth
- [x] Integrate Auth0 authentication
- [x] Correctly validate auth token
- [x] Automate refresh token usage
- [ ] Auth-Store
  - [x] Enable local storage of auth tokens
  - [ ] Ensure store is safely encrypted
  
### UI
- [ ] Header
  - [ ] Add logo
  - [x] Add user dropdown menu
  - [ ] Add links for various games
    - [ ] DDR
    - [ ] DRS
- [ ] Footer
  - [x] API status
  - [ ] Licensing
  - [ ] Cookie notice
- [ ] Decide on general style approach

### User
- [ ] Provide eagate API integration
  - [x] Retrieve eagate link status
  - [ ] EaAccount Linking
    - [x] Allow linking BST profile to Eagate profile
    - [ ] Allow linking BST profile to multiple Eagate profiles
  - [x] Allow unlinking eagate profiles
  - [ ] Improve user page performance
    - [ ] Async ea link status
  - [ ] Improve linking ux
- [ ] Profile automation
  - [ ] Opt-In third-party update
  - [ ] Automatic update user

### DDR
- [ ] Provide ddr API integration
  - [ ] Profile refresh
  - [x] Profile update
  - [ ] Song retrieval
  - [ ] Score retrieval
- [ ] Base ddr page
  - [ ] User details
  - [ ] Profile update
  - [ ] Profile refresh
- [ ] Provide song overview
  - [ ] Minified song jacket
  - [ ] Score overview
  - [ ] Sorting
- [ ] Provide song page
  - [ ] Chart overview
  - [ ] Scores for linked accounts
  - [ ] Song jacket

### General
- [ ] Improve readme
- [ ] Add dev build task
- [x] Add prod build task
- [ ] Redeploy to a more performant server
  - [ ] Decide on best region
- [x] Automatic cert

### Refactoring
- [ ] Identify code optimisations
- [ ] Improve templating