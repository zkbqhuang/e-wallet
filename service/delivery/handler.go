package delivery

import (
	"e-wallet/service/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/securecookie"
)

const (
	// tampilan untuk idnex (pada saat endpoint / diakses)
	// hal ini dibuat agar lebih mudah jika ingin mengakses melalui web agar lebih mudah
	indexPage = `
	<h1>Login</h1>
	<form method="post" action="/login">
		<label for="username">User name</label>
		<input type="text" id="username" name="username">
		<label for="password">Password</label>
		<input type="password" id="password" name="password">
		<button type="submit">Login</button>
	</form>
	`

	// tampilan web untuk menampilkan username setelah berhasil login
	// hal ini dibuat agar lebih mudah jika ingin mengakses melalui web agar lebih mudah
	internalPage = `
	<h1>Internal</h1>
	<hr>
	<small>User: %s</small>
	<form method="post" action="/logout">
		<button type="submit">Logout</button>
	</form>
	`
)

// cookie handling
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["username"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"username": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

// deklarasi dan handling home ("/")
func homeHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

// deklarasi dan handling untuk success login ("/internal")
func pageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

// deklarasi dan handling login ("/login")
func loginHandler(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if username != "" && pass != "" {
		// .. check credentials ..
		setSession(username, response)
		redirectTarget = "/internal"
	}
	http.Redirect(response, request, redirectTarget, 302)
}

// deklarasi dan handling login ("/logout")
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

// top-up handler
func topupHandler(response http.ResponseWriter, request *http.Request) {
	usernamex := getUserName(request)

	username := request.FormValue("username")
	balance := request.FormValue("jumlah")
	var status string

	if usernamex != "" {
		jumlah, _ := strconv.Atoi(balance)
		err := repository.TopupBalanceRepository(username, jumlah)
		if err != nil {
			println("Error while exec")
			println(err.Error())
			status = "Failed"
		} else {
			status = "Success"
		}
	} else {
		status = "Before we go, please login first"
	}

	fmt.Fprintf(response, status)
}
