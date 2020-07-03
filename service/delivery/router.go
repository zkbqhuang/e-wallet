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
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/internal", pageHandler)

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	router.HandleFunc("/top-up", topupHandler).Methods("POST")

	http.Handle("/", router)
	fmt.Println("server started at localhost:8000")
	http.ListenAndServe(":8000", nil)
}
