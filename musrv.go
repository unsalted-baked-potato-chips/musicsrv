package main

import (
    "fmt"
    "log"
    "net/http"
)

func searchForSong(s string){

}

func handlePost(w *http.ResponseWriter, r *http.Request){
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(*w, "ParseForm() err: %v", err)
		return
	}
    switch r.FormValue("playerControl") {
        case "play":
            fmt.Fprintf(*w, "Playing")
        case "pause":
            fmt.Fprintf(*w, "Pausing")
        case "skip":
            fmt.Fprintf(*w, "Skipping")
        default "":
        http.Error(*w, "Unsupported Player Control", http.StatusNotFound)
    }
}

func requestHandler(w http.ResponseWriter, r *http.Request){
    if r.URL.Path != "/" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }

    switch r.Method {
        case "GET":
            q := r.URL.Query().Get("search")
            if q == ""{
                http.ServeFile(w, r, "index.html")
                return
            }
            fmt.Fprintf(w, q)
            searchForSong(q)
        case "POST":
            handlePost(&w, r)
        default:
            http.Error(w, "Method is not supported.", http.StatusNotFound)
    }
}

func main () {
    http.HandleFunc("/", requestHandler)
    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
