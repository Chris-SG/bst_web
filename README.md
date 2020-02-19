**BST Web**

WIP web-server leveraging Auth0 or another identity provider.

Build for *nix with:

```
env GOOS=linux GOARCH=amd64 go build
``` 

Run with:

```
./bst_web -static="./dist" -entry="./dist/index.html" -clientid="authclientid" -clientsecret="authclientsecret" -issuer="https://authissuer.com/" -host="abc.com"
```