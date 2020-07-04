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
)

// cookie handling
var cookieHandler = securecookie.New(
	// perintah untuk membuat key secara random
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// perintah untuk handling atau fungsi pengecekan after login
// setelah login pasti data username akan tersimpan di cookie
// jika fungsi ini dipanggil tidak mengembailkan dara username maka user blm berhasil login
func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		// deklarasi map untuk menampung value dari cookie yang akan disimpan
		// deklarasi ini berbentuk map dengan key bertipe string dan value nya bertipe string
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["username"]
		}
	}
	// nilai balikan pada saat fungsi ini dipanggil
	// value yang dibalikan adalah data username yang bertipe string
	return userName
}

// deklarasi handling untuk meng-set sesion pada saat login
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

// handling untuk menghandle logout, dimana proses nya adalah menghapus sesion yang tersimpan dengan data login tersebut
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

// deklarasi dan handling login ("/login")
func loginHandler(response http.ResponseWriter, request *http.Request) {
	var status string

	username := request.FormValue("username")
	pass := request.FormValue("password")
	userData, err := repository.GetUserData(username)
	if err != nil {
		println(err.Error())
	}

	// alasan saya tidak melakukan enkripsi pada password karena untuk project sekecil ini dan se simple ini tidak diperlukan untuk enkripsi dan karena keterbatas waktu juga
	// jika nantinya dalam project asli pasti saya akan melakukan enkripsi
	if userData.Username == username && userData.Password == pass {
		if username != "" && pass != "" {
			// perintah untuk melakukan pengecekan credential
			// jika credential nya benar maka akan diarahkan ke alamat "/internal"
			setSession(username, response)
			status = "Berhasil login"
		}
	} else {
		status = "Maaf, kombinasi yang anda masukan salah"
	}
	fmt.Fprintf(response, status)
}

// deklarasi dan handling login ("/logout")
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	status := "Berhasil logout"
	fmt.Fprintf(response, status)
}

// top-up handler untuk menghandle route "/top-up"
func topupHandler(response http.ResponseWriter, request *http.Request) {
	// perintah untuk melakukan pengecekan credential
	// jika fungsi ini dipanggil dan mengembalikan value yang tidak sama atau kosong maka variable ini akan dijadikan parameter pengecekan
	// jika data variable "usernamex" kosong atau tidak sesuai maka tidak dapat melakukan action utama yaitu action top-up
	usernamex := getUserName(request)

	// mengambil nilai dengan key yang tertera dibawah, data yang ditangkap adalah data yang dimasukan sebagai paramter pada saat memanggil route
	username := request.FormValue("username")
	balance := request.FormValue("jumlah")
	var status string

	// pengecekan credential, jika data sesuai maka akan dapat meng-eksekusi fungsi top-up
	if usernamex != "" {
		// convert value jumlah dari inputan berupa string dirubah menjadi int
		jumlah, _ := strconv.Atoi(balance)
		// pemanggilan fungsi pada repository untuk melakukan eksekusi yang akan dijalankan pada repository, yaitu berhubungan dengan database
		err := repository.TopupBalanceRepository(username, jumlah)
		if err != nil {
			println("Error while exec")
			println(err.Error())
			status = "Failed"
		} else {
			status = "Success"
		}
	} else { // statment jika syarat credential diatas tidak terpenuhi
		status = "Before we go, please login first"
	}

	fmt.Fprintf(response, status)
}

// transfer handler
func transferHandler(response http.ResponseWriter, request *http.Request) {
	// perintah untuk melakukan pengecekan credential
	// jika fungsi ini dipanggil dan mengembalikan value yang tidak sama atau kosong maka variable ini akan dijadikan parameter pengecekan
	// jika data variable "usernamex" kosong atau tidak sesuai maka tidak dapat melakukan action utama yaitu action transfer
	usernamex := getUserName(request)

	// mengambil nilai dengan key yang tertera dibawah, data yang ditangkap adalah data yang dimasukan sebagai paramter pada saat memanggil route
	username := request.FormValue("username")
	usernameTujuan := request.FormValue("tujuan")
	balance := request.FormValue("jumlah")

	// deklarasi variable untuk menerima nilai balikan / return value dari fungsi transfer pada scope repository
	var status string

	// pengecekan credential, jika data sesuai maka akan dapat meng-eksekusi fungsi transfer
	if usernamex != "" {
		jumlah, _ := strconv.Atoi(balance)
		// pemanggilan fungsi pada repository untuk melakukan eksekusi yang akan dijalankan pada repository, yaitu berhubungan dengan database
		err := repository.TransferBalanceRepository(username, usernameTujuan, jumlah)
		if err != nil {
			println("Error while exec")
			println(err.Error())
			status = "Failed"
		} else {
			status = "Success"
		}
	} else { // statment jika syarat credential diatas tidak terpenuhi
		status = "Before we go, please login first"
	}

	fmt.Fprintf(response, status)
}
