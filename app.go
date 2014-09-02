package main
import (
	"fmt"
	"net/http"
	"encoding/json"
    "io/ioutil"
    "strconv"
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
    top_channel = time.NewTicker(time.Second).C
    vote_channel = make(chan int,1000)
    go process_vote()
    go sync_to_file()
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8088", nil)
}
