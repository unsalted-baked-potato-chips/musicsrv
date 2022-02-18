package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os/exec"
    "io/ioutil"
    "bytes"
    "html/template"
    "encoding/json"
)
type YTID struct {
    Kind string `json:"kind"`
    VideoID string `json:"videoId"`
}
type YTSnippet struct {
    Title string `json:"title"`
}
type YTResult struct {
    Snippet YTSnippet `json:"snippet"`
    ID YTID `json:"id"`
}
type YTResponse struct {
    Results []YTResult `json:"items"`
}

var ytkey string

func enqueueSong(s string) string{
    cmd := exec.Command("youtube-dl", s, "--get-filename")
    var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
    return ""
    cmd = exec.Command("umpv", "add",  out.String());
    cmd.Stdout = &out
	err = cmd.Run()


	if err != nil {
		log.Println(err)
	}

    return out.String()
}

func searchYT(s string, w *http.ResponseWriter){
    u, err := url.Parse("https://www.googleapis.com/youtube/v3/search")
    if err != nil {
        return

    }
    // Query params
    params := url.Values{}
    params.Add("q",s)
    params.Add("type","video")
    params.Add("part","snippet")
    params.Add("key",ytkey)
    
    u.RawQuery = params.Encode()

    r, err := http.Get(u.String())
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        log.Fatalln(err)
    }
    

    var ytr YTResponse
    err = json.Unmarshal(body, &ytr)
    if err != nil {
        fmt.Fprintf(*w, "%s", body)
        log.Fatal(err)
    }

    t := template.Must(template.ParseFiles("results.html"))
    t.Execute(*w, ytr)
}

func handlePlayerButtons(w *http.ResponseWriter, r *http.Request){
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
        default :
        http.Error(*w, "Unsupported Player Control", http.StatusNotFound)
    }
}

func requestHandler(w http.ResponseWriter, r *http.Request){

    switch r.Method {
        case "GET":
            q := r.URL.Query().Get("search")
            if q == ""{
                http.ServeFile(w, r, "index.html")
                return
            }
            searchYT(q, &w)
        case "POST":
            switch r.URL.Path{
            case "/":
                handlePlayerButtons(&w,r)
            case "/enqueue/":
               fmt.Fprintf(w, "%s", r.FormValue("vidURL"))
            default:
                http.Error(w, "404 Not Found.", http.StatusNotFound)
                return
            }
        default:
            http.Error(w, "Method is not supported.", http.StatusNotFound)
    }
}

func main () {
    ytkb, err := ioutil.ReadFile("apikey")
    if err != nil {
        log.Fatal(err)

    }

    // Convert []byte to string and print to screen
    ytkey= string(ytkb)
    http.HandleFunc("/", requestHandler)
    fmt.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}
