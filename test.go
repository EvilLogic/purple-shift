package main

import (
  "context"
  "net"
  "strings"
  "sync"
  "time"
  "log"
  "net/http"
  "io/ioutil"
  "fmt"
  "encoding/json"
  
  "golang.org/x/sync/semaphore"
)

type Team struct {
	Name		string	`json:"name"`
	URL			string 	`json:"url"`
	Number		int		`json:"number"`
	BPoints		int		`json:"bpoints"`
	RPoints		int		`json:"rpoints"`
	IScore		int		`json:"iscore"`
	Total		int		`json:"total"`
	//Scans		[]Scan	`json:"scans"`
}

type Inject struct {
	Title		string	`json:"title"`
	Question	string	`json:"question"`
	Points		int		`json:"points"`
	Answer		string	`json:"answer"`
}

type Scan struct {
	Name		string	`json:"name"`
	Port		int		`json:"port"`
	Points		int		`json:"points"`
	Interval	int		`json:"interval"`
}

type Answer struct {
	Title 		string	`json:"title"`
	Answer		string	`json:"answer"`
	Team		string	`json:"team"`
}

var teams 		[]Team
var injects 	[]Inject
var scans 		[]Scan
var interval	int

func main() {
	http.HandleFunc("/css/", serveCSS)
	http.HandleFunc("/js/", serveJS)
	http.HandleFunc("/inject/", injectHandler)
	http.HandleFunc("/api/info/", infoHandler)
	http.HandleFunc("/scan/", scanHandler)
	http.HandleFunc("/", viewHandler)
	mapinit()
	log.Fatal(http.ListenAndServe(":8080", nil))
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
    p, _ := loadPage(title)
    fmt.Fprintf(w, "%s", p)
}

func serveCSS(w http.ResponseWriter, r *http.Request){
    w.Header().Set("Content-Type", "text/css")

    title := r.URL.Path[1:]
    p, _ := loadPage(title)
    fmt.Fprintf(w, "%s", p)
}

func serveJS(w http.ResponseWriter, r *http.Request){
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
		if e.Title == answer.Title && e.Answer == answer.Answer{
			for i, _ := range teams {
				if teams[i].Name == answer.Team {
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

func mapinit(){
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
	teams[0].BPoints = 1000
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

func scanTeams() {
	ps := &PortScanner{
		tlist:  teams,
		lock: 	semaphore.NewWeighted(512),
	}
	ps.Start(500*time.Millisecond)
}

type PortScanner struct {
	tlist  []Team
	lock 	*semaphore.Weighted
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
			go func(target string) {
				defer ps.lock.Release(1)
				defer wg.Done()
				
				if ScanPort(target, timeout) == "open" {
					log.Print(target)
					teams[i].BPoints += s.Points
				} else {
					teams[i].RPoints += s.Points
				}
			}(target)
		}
	}
}
