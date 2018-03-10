package main

import (
	"net/http"
	"github.com/codegangsta/negroni"
)


var holder *Holder
func main()  {
	holder = newHolder()
	http.Handle("/c", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))
	http.ListenAndServe(":8888", nil)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Origin")
	serverWs(holder,w,r)
}

