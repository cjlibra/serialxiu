package main 
import ( 
	"net/http" 
	"log" 
	"fmt" 
	"runtime" 
	"io/ioutil" 
	"text/template"
) 
func handler514(w http.ResponseWriter, r *http.Request) {

     t := template.New("struct data demo template") 
	 t, _ = t.Parse("<h1>hello, {{.UserName}}! </h1>") 
	 actorMap := make(map[string]string)
	 actorMap["UserName"]="sssssss"
	 t.Execute(w, actorMap)


}
func handler513(w http.ResponseWriter, r *http.Request) { 
	w.Header().Set("Connection", "keep-alive") 
	//a := []byte("aaaaa...512...aaa") 
	b,_ := ioutil.ReadFile("222.bmp")
	//w.Header().Set("Content-Length", fmt.Sprintf("%d", len(a))) 
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b))) 
	w.Header().Set("Content-Type", "text/html") 
	//w.Write(a) 
	w.Write(b) 
} 

func main() { 
	runtime.GOMAXPROCS(runtime.NumCPU()) 
	fmt.Println("cpu is",runtime.NumCPU())
	http.Handle("/", http.FileServer(http.Dir("tmp")))
	 
    http.HandleFunc("/514", handler514) 
	http.HandleFunc("/512b", handler513) 
	log.Fatal(http.ListenAndServe(":8080", nil)) 
}