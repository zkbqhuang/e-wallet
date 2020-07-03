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

	http.Handle("/", router)
	fmt.Println("server started at localhost:8000")
	http.ListenAndServe(":8000", nil)
}
