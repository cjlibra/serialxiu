package main

import (
     "fmt"
	 "path/filepath"
	 "os"
	 "io/ioutil"
	// "path"
)



func main(){

   // filepath.Walk(".\\", we)
   dir, _ := ioutil.ReadDir(".\\")
   for _ ,file := range dir {
      file1,_:=filepath.Abs(file.Name())
	  //file1.IsDir()
      fmt.Println(file1,file.IsDir()?"mulu":"wenjian")
	
   
   }

}


func  we(path string, fi os.FileInfo, err error) error {
	if nil == fi {
	return err
	}
	if fi.IsDir() {
	return nil
	}
    name := fi.Name()
	fmt.Println(path +"\\"+name)
	

    return nil
}