package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

// CreateBstApiRouter will generate a router mapped against BST API. Middleware
// may be passed in to then be used by certain routes.
func CreateAjaxRouter(prefix string, middleware map[string]*negroni.Negroni) *mux.Router {
	ajaxRouter := mux.NewRouter().PathPrefix(prefix + "/ajax").Subrouter()
	ajaxRouter.Path("/eagate_login_status").Handler(negroni.New(
		negroni.Wrap(http.HandlerFunc(EagateLoginStatusGet)))).Methods(http.MethodGet)

	return ajaxRouter
}

func EagateLoginStatusGet(rw http.ResponseWriter, r *http.Request) {
	token, err := TokenForRequest(r)
	if err != nil {
		text := `<a id="eagate-login-state">Please ensure you are logged in to BST.</a>`
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(text))
		return
	}
	status, users := EagateLoginGetImpl(token)

	if status.Status == "bad" || len(users) == 0 {
		text := `
		<a id="eagate-login-state">
        <form class="form-inline">
            <div class="form-row align-items-center">
                <input id="eagate-username" class="form-control form-control-sm" type="text" placeholder="Eagate Username">
                <input id="eagate-password" class="form-control form-control-sm" type="password" placeholder="Eagate Password">
                <button class="btn btn-primary" type="submit" onClick="eagateLogin()">Login</button>
            </div>
        </form>
		</a>`
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(text))
		return
	}

	var text []byte
	text = append(text, []byte(`<a id="eagate-login-state">`)...)
	for _, user := range users {
		userText := fmt.Sprintf(`
        <p>Currently linked to %s.</p>
        <button type="button" class="btn btn-primary" onClick="eagateLogout('%s')">Unlink</button>`, user.Username, user.Username)
		text = append(text, []byte(userText)...)
	}
	text = append(text, []byte(`</a>`)...)
	rw.WriteHeader(http.StatusOK)
	rw.Write(text)
	return
}