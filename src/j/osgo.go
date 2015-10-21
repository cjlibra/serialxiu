package main
import (
   "fmt"
    "os"
   // "reflect"
   //"error"
)


func main() {

    f, err := os.Open("1.go")
	err = f.Close()
	if err != nil {
	   fmt.Println("err",err.Error())
	}

    defer f.Close()

}
