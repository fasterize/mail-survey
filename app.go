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
var top_channel <- chan time.Time
var vote_channel chan int
func process_vote(){
    for ;; {
    id :=<- vote_channel
    conf.Resultat[id].Votes+=1
    fmt.Println("new vote ",id)
    }
}
func sync_to_file(){
    for{
        <-top_channel
        str,err := json.Marshal(conf)
        if(err== nil){
        fmt.Println("hello",conf,string(str))
        ioutil.WriteFile("compagne.json",str,0644)
        } else {
            fmt.Println("error sync_to_file ",err)
        }
    }

}

func handler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[1:]
    fmt.Println(id)
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
    i,err := strconv.Atoi(id)
    if(err == nil){
        vote_channel <- i
    }


}
func getResult(w http.ResponseWriter, r *http.Request) {
    content, err := ioutil.ReadFile("compagne.json")
    if err!=nil{
        fmt.Print("Error:",err)

    }
    fmt.Fprintf(w, string(content))


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
    top_channel = time.NewTicker(time.Second*30).C
    vote_channel = make(chan int,1000)
    go process_vote()
    go sync_to_file()
    http.HandleFunc("/status", getResult)
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
