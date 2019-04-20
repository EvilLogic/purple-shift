package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

type Team struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Number  int      `json:"number"`
	BPoints int      `json:"bpoints"`
	RPoints int      `json:"rpoints"`
	Injects []string `json:"injects"`
	IScore  int      `json:"iscore"`
	Total   int      `json:"total"`
	//Scans		[]Scan		`json:"scans"`
}

type Inject struct {
	Title    string `json:"title"`
	Question string `json:"question"`
	Points   int    `json:"points"`
	Answer   string `json:"answer"`
}

type Scan struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Points   int    `json:"points"`
	Interval int    `json:"interval"`
}

type Answer struct {
	Title  string `json:"title"`
	Answer string `json:"answer"`
	Team   string `json:"team"`
}

var teams []Team
var injects []Inject
var scans []Scan
var interval int
var ticker *time.Ticker

func main() {
	http.HandleFunc("/css/", serveCSS)
	http.HandleFunc("/js/", serveJS)
	http.HandleFunc("/settings/", settingsHandler)
	http.HandleFunc("/inject/", injectHandler)
	http.HandleFunc("/api/info/", infoHandler)
	http.HandleFunc("/scan/", scanHandler)
	http.HandleFunc("/", viewHandler)
	mapinit()
	//startScans()
	log.Fatal(http.ListenAndServe(":80", nil))
}

func loadPage(title string) ([]byte, error) {
	filename := title
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[1:]

	if title == "" {
		title = "test.html"
	}

	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p)
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")

	title := r.URL.Path[1:]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p)
}

func serveJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")

	title := r.URL.Path[1:]
	p, _ := loadPage(title)
	fmt.Fprintf(w, "%s", p)
}

func injectHandler(w http.ResponseWriter, r *http.Request) {
	var answer Answer
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &answer)

	//title := r.URL.Path[8:]
	//body, _ := ioutil.ReadAll(r.Body)
	//answer := string(body)
	log.Print(answer)

	for _, e := range injects {
		if e.Title == answer.Title && strings.Contains(strings.ToLower(answer.Answer), strings.ToLower(e.Answer)) {
			for i, _ := range teams {
				if teams[i].Name == answer.Team {
					if contains(e.Title, teams[i]) {
						fmt.Fprintf(w, "already solved")
						return
					}
					teams[i].Injects = append(teams[i].Injects, e.Title)
					teams[i].IScore += e.Points
					fmt.Fprintf(w, "correct")
					return
				}
			}
			fmt.Fprintf(w, "invalid team")
			return
		}
	}
	fmt.Fprintf(w, "incorrect")
}

func contains(ti string, te Team) bool {
	for _, e := range te.Injects {
		if e == ti {
			return true
		}
	}
	return false
}

func mapinit() {
	t, _ := loadPage("api/teams.json")
	json.Unmarshal(t, &teams)
	t, _ = loadPage("api/injects.json")
	json.Unmarshal(t, &injects)
	t, _ = loadPage("api/scans.json")
	json.Unmarshal(t, &scans)

	//for i, _ := range teams{
	//	teams[i].Scans = scans
	//}

	fmt.Println("%#v", teams)
	fmt.Println("%#v", injects)
}

func getJSON(v interface{}) string {
	blob, _ := json.Marshal(v)
	return string(blob)
}

func getSecret() []Inject {
	secret := make([]Inject, len(injects))
	copy(secret, injects)

	for i, _ := range secret {
		secret[i].Answer = "You silly hacker"
	}
	return secret
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	pw := r.Header.Get("Password")
	body, _ := ioutil.ReadAll(r.Body)

	if pw != "BuildYourOwnDamnBirdfeeder" {
		fmt.Fprintf(w, "MESS WITH THE BEST DIE LIKE THE REST")
		return
	}

	if len(body) > 0 {
		path := r.URL.Path[10:]

		if path == "teams" {
			json.Unmarshal(body, &teams)
			fmt.Fprintf(w, getJSON(teams))
		} else if path == "injects" {
			json.Unmarshal(body, &injects)
			fmt.Fprintf(w, getJSON(injects))
		} else if path == "scans" {
			json.Unmarshal(body, &scans)
			fmt.Fprintf(w, getJSON(scans))
		}
	} else {
		path := r.URL.Path[10:]

		if path == "teams" {
			fmt.Fprintf(w, getJSON(teams))
		} else if path == "injects" {
			fmt.Fprintf(w, getJSON(injects))
		} else if path == "scans" {
			fmt.Fprintf(w, getJSON(scans))
		}
	}
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[10:]
	//p, _ := loadPage(title)
	if path == "teams" {
		fmt.Fprintf(w, getJSON(teams))
	} else if path == "injects" {
		secret := getSecret()
		fmt.Fprintf(w, getJSON(secret))
	} else if path == "scans" {
		fmt.Fprintf(w, getJSON(scans))
	}
}

func scanHandler(w http.ResponseWriter, r *http.Request) {
	scanTeams()
	fmt.Fprintf(w, "it works!")
}

func startScans() {
	ticker = time.NewTicker(30 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				scanTeams()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func scanTeams() {
	ps := &PortScanner{
		tlist: teams,
		lock:  semaphore.NewWeighted(512),
	}
	ps.Start(3000 * time.Millisecond)
}

type PortScanner struct {
	tlist []Team
	lock  *semaphore.Weighted
}

func ScanPort(target string, timeout time.Duration) string {
	var result string

	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			log.Print("too many files!")
			time.Sleep(timeout)
			ScanPort(target, timeout)
		} else {
			result = "closed"
		}
		return result
	}

	conn.Close()
	return "open"
}

func (ps *PortScanner) Start(timeout time.Duration) {
	wg := sync.WaitGroup{}
	defer wg.Wait()

	for i, t := range ps.tlist {
		for _, s := range scans {
			target := fmt.Sprintf("%s:%d", t.URL, s.Port)

			ps.lock.Acquire(context.TODO(), 1)
			wg.Add(1)
			go func(target string, team *Team) {
				defer ps.lock.Release(1)
				defer wg.Done()

				if ScanPort(target, timeout) == "open" {
					log.Print(target + " open")
					team.BPoints += s.Points
				} else {
					log.Print(target + " close")
					team.RPoints += s.Points
				}
			}(target, &teams[i])
		}
	}
}
