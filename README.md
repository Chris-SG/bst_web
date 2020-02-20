**BST Web**

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
    -callback="/callback" \
    -filestorekey="averysecurekey" \
    -servehost="my.host.com" \
    -serveport="443"
```

---

*TODO*
 - Improve this readme
 - Rewrite routing
 - Add a lot of pages
 - Connect to bst-api
 - ???