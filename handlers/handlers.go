package handlers

import (
	"fmt"
	"github.com/bpina/shortened/data"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
	"strconv"
)

func ShortenHandler(w http.ResponseWriter, req *http.Request) {
	originalUrl := req.URL.Query().Get("url")

	originalUrl, err := url.QueryUnescape(originalUrl)

	if err != nil {
		SendJsonError(w, "could not parse source url")
	}

	id, err := data.StoreUrl(originalUrl)

	if err != nil {
		SendJsonError(w, "failed to create id")
		return
	}

	encodedId := strconv.FormatInt(id, 36)

	response := fmt.Sprintf("{status: 200, id: \"%s\"}", encodedId)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

func UrlRedirectHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	if id == "" {
		SendTextError(w, "could not determine id")
		return
	}

	decodedId, err := strconv.ParseInt(id, 36, 64)

	if err != nil {
		SendTextError(w, "could not parse id")
		return
	}

	originalUrl, err := data.GetUrl(decodedId)

	if err != nil {
		SendTextError(w, err.Error())
		return
	}

	w.Header().Set("Location", originalUrl)
	w.WriteHeader(http.StatusFound)
	w.Write([]byte{})

	go func(id int64) {
		_, err = data.UpdateUrlUsage(id)
		if err != nil {
			fmt.Println(err)
		}
	}(decodedId)

}

func SendJsonError(w http.ResponseWriter, message string) {
	fmt.Println(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	error := fmt.Sprintf("{status: 500, error: \"%s\"}", message)
	w.Write([]byte(error))
}

func SendTextError(w http.ResponseWriter, message string) {
	fmt.Println(message)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
