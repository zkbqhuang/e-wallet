package delivery

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// inisiasi router package
var router = mux.NewRouter()

// Router init
func Router() {
	// route untuk home
	router.HandleFunc("/", homeHandler)
	// route untuk status setelah berhasil login
	router.HandleFunc("/internal", pageHandler)

	// route untuk melakukan login / create cookie & sesion
	router.HandleFunc("/login", loginHandler).Methods("POST")
	// route untuk melakukan logout / menghapus sesion yang ada sesuai dengan credential login yang digunakan
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	// route untuk melakukan top-up
	router.HandleFunc("/top-up", topupHandler).Methods("POST")
	// route untuk melakukan transfer
	router.HandleFunc("/transfer", transferHandler).Methods("POST")

	// base route
	http.Handle("/", router)
	// status yang menandakan bahwa program ini sudah berjalan pada detail seperti dibawah
	fmt.Println("server started at localhost:8000")
	// define port yang digunakan, disini menggunakan port :8000, bisa diubah sesuai keinginan
	http.ListenAndServe(":8000", nil)
}
