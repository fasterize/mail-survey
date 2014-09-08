package main
import (
	"fmt"
	"net/http"
	"encoding/json"
    "io/ioutil"
    "strconv"
    "os"
    "time"
)

type Stat struct {
	Votes int
}

type All struct {
	Resultat [10]Stat
}
var conf All
var sync_channel chan bool
var timer_channel <- chan time.Time

var vote_channel chan int
func process_vote(){

    for ;; {
    id :=<- vote_channel
    if id>=0 && id <10 {
        conf.Resultat[id].Votes+=1
        fmt.Println("new vote ",id)
        }
    }
}
func sync_to_file(){
    for{
        select {
            case <-timer_channel:
                save_to_file()
            case <-sync_channel:
                save_to_file()
            }
    }

}
func save_to_file(){
    fmt.Println("save to file")
    str,err := json.Marshal(conf)
        if(err== nil){
        fmt.Println("hello",conf,string(str))
        ioutil.WriteFile("compagne.json",str,0644)
        } else {
            fmt.Println("error sync_to_file ",err)
        }
}

func handler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[1:]
    http.Redirect(w,r,"http://www.fasterize.com/en/survey_thanks",301)
    i,err := strconv.Atoi(id)
    if(err == nil){
        vote_channel <- i
    }


}
func getResult(w http.ResponseWriter, r *http.Request) {
    response := "<h1> Vote result</h1>"
    for i:=0;i<10;i++{
        response+= fmt.Sprintf("<span style=\"color:red;margin-right:20px\">Choix %d:</span><span style=\"color:green\">%d hits</span></br>",i,conf.Resultat[i].Votes)
    }
    fmt.Fprintf(w, response)


}

func initCounter(w http.ResponseWriter, r *http.Request) {
   for i:=0;i<10;i++{
       conf.Resultat[i].Votes=0
    }
    fmt.Println("init counter")
    sync_channel <- true
    http.Redirect(w,r,"/status",301)
}

func main() {
    content, err := ioutil.ReadFile("compagne.json")
    if err!=nil{
        fmt.Print("Error:",err)

    } else {
      err=json.Unmarshal(content, &conf)
      if err!=nil{
            fmt.Print("Error:",err)
        }
    }
    timer_channel = time.NewTicker(time.Second*30).C
    vote_channel = make(chan int,1000)
    sync_channel = make(chan bool)
    go process_vote()
    go sync_to_file()
    http.HandleFunc("/status", getResult)
    http.HandleFunc("/initCounterFF00123", initCounter)
    http.HandleFunc("/", handler)
    http.ListenAndServe(GetPort(), nil)
}
func GetPort() string {
        var port = os.Getenv("PORT")
        // Set a default port if there is nothing in the environment
        if port == "" {
                port = "4747"
                fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
        }
        return ":" + port
}
