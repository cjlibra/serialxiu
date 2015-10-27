package main

import (
	"os"
	"encoding/json"
	"net/http"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"   

	  "github.com/tarm/goserial"                             
)

var Input_read = flag.String("r", "020304F", "input reader id")
var Input_write = flag.String("w", "020304F", "input writer id")
var Input_file = flag.String("f", "filename", "input filename")
var Input_scan = flag.String("s", "FFF", "input id for scan")
var Input_daemon = flag.String("d", "ddd", "no input")

var senddata []byte
var srwflag = 0
type OUTSTRC struct{
	Errorint int
	Errstr string
	Rfid []string
	
	
}


func jreadrf(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	
		 
	readid := r.FormValue("readid")
	com := r.FormValue("com")
	rfid := r.FormValue("rfid")
	fmt.Println("read tag data")
	
	var outstr OUTSTRC
	outstr.Errorint = 1
	outstr.Errstr = "串口打不开"
	b, _ := json.Marshal(outstr)
	fmt.Println(com," ",readid)
	c := &serial.Config{Name: com, Baud: 38400}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println("error ",com)
		w.Write([]byte(b))
		return
	}
	
	
	senddata = []byte("\x00\x02\x09\x02\x03\x04\x00\x01\x00")
	onechars, _ := hex.DecodeString(readid)
	senddata[0] = onechars[0]
	onechars, _ = hex.DecodeString(rfid)
	senddata[3] = onechars[0]
	senddata[4] = onechars[1]
	senddata[5] = onechars[2]
	senddata[8] = Xor(string(senddata[0:8]))
	 
	filecontent := make([]byte, 16*64)
	slicecount := 0
	for i0 := 0; i0 < 64; i0++ {
		senddata[7] = byte(i0)
		fmt.Println(senddata[7])
		senddata[8] = Xor(string(senddata[0:8]))
		n, err := s.Write(senddata)

		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 128)
		recvdata := make([]byte, 128)
		num := 0

		for {
			n, err = s.Read(buf)

			for i1 := 0; i1 < n; i1++ {
				recvdata[i1+num] = buf[i1]
			}
			num = num + n
			if num >= 2 {
				intnum := int(recvdata[1])
				if num == intnum {
					fmt.Println(recvdata)
					if readcomppack(recvdata) == 0 || recvdata[num-1] != Xor(string(recvdata[:num-1])) {
						i0 = i0 - 1
						fmt.Println(recvdata[num-1])
						fmt.Println(Xor(string(recvdata[:num-1])))
					} else {
						for i2 := 0; i2 < 16; i2++ {
							filecontent[i2+16*slicecount] = recvdata[i2+4]
						}
						slicecount++

					}

					break
				}
			}

		}

	}
	s.Close()
	currentDir,_ := os.Getwd()
	writefilelength, _ := strconv.Atoi(string(filecontent[0:4]))
	//fmt.Println(writefilelength)
	if filecontent[4+writefilelength] == Xor(string(filecontent[0:4+writefilelength])) {
		ioutil.WriteFile("outputFile"+rfid, filecontent[4:writefilelength+4], 0x644)
		outstr.Errorint = 0;
		outstr.Errstr = "文件生成成功，文件名为：" +currentDir+ "\\outputFile"+rfid;
		b, _ := json.Marshal(outstr)
		w.Write([]byte(b))
	} else {
		fmt.Println("file checksum error")
		outstr.Errorint = 1;
		outstr.Errstr = "file checksum error";
		b, _ := json.Marshal(outstr)
		w.Write([]byte(b))
		
	}
	
}
func jwriterf(w http.ResponseWriter, r *http.Request) {
	
	r.ParseForm() //解析参数，默认是不会解析的
		 
	readid := r.FormValue("readid")
	com := r.FormValue("com")
	rfid := r.FormValue("rfid")
	fmt.Println("aaaaaaaa")
	file, _, err := r.FormFile("file")
	fmt.Println("aaaaaaaab",com)
	
	
	if err != nil {
		fmt.Println("file upload error")
		return
	}
	defer file.Close()
	var outstr OUTSTRC
	outstr.Errorint = 1
	outstr.Errstr = "串口打不开"
	outstr.Rfid=[]string{"",""}
	b, _ := json.Marshal(outstr)
	fmt.Println(com," ",readid," ",rfid)
	c := &serial.Config{Name: com, Baud: 38400}
	s, err := serial.OpenPort(c)
	defer s.Close()
	if err != nil {
		fmt.Println("error ",com,outstr)
		w.Write([]byte(b))
		return
	}
	if rfid == "" {
	  outstr.Errorint = 1
	  outstr.Errstr = "rfid无数据"
	  outstr.Rfid=[]string{"",""}
	  b, _ = json.Marshal(outstr)
	   fmt.Println("error ",com,outstr)
	  w.Write([]byte(b))
	  return
	}
	
	fmt.Println("write tag data")
	senddata = make([]byte, 25)
	onechars, _ := hex.DecodeString(readid)
	senddata[0] = onechars[0]
	senddata[1] = '\x03'
	senddata[2] = '\x19'
	onechars, _ = hex.DecodeString(rfid)
	senddata[3] = onechars[0]
	senddata[4] = onechars[1]
	senddata[5] = onechars[2]

	 
	filecontent, err := ioutil.ReadAll(file)
	filelength := len(filecontent)
	if err != nil {
		fmt.Println("file can not read")
		return
	}
	filecontent1024 := make([]byte, 1024)
	copy(filecontent1024[0:4], []byte(fmt.Sprintf("%04d", filelength)))
	for aa := 0; aa < len(filecontent); aa++ {
		filecontent1024[aa+4] = filecontent[aa]
	}
	filecontent1024[4+filelength] = Xor(string(filecontent1024[:4+filelength]))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(filelength)
	for slicecount := 0; slicecount < 64; slicecount++ {

		senddata[7] = byte(slicecount)

		for ii := 0; ii < 16; ii++ {
			senddata[ii+8] = filecontent1024[slicecount*16+ii]
		}
		senddata[24] = Xor(string(senddata[0:24]))
		fmt.Println(senddata)
	sendagain:

		n, err := s.Write(senddata)

		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 128)
		recvdata := make([]byte, 128)
		num := 0
		for {
			n, err = s.Read(buf)
			for i := 0; i < n; i++ {
				recvdata[i+num] = buf[i]
			}
			num = num + n
			if num >= 2 {
				intnum := int(recvdata[1])
				if num == intnum {
					break
				}
			}

		}
		if recvdata[2] == '\x01' {
			fmt.Println("...")
			time.Sleep(time.Second * 1)
			goto sendagain
		}
		if recvdata[2] == '\x00' {
			fmt.Println("ok", recvdata[7])
			if slicecount == 63 {
				outstr.Errorint = 0
				outstr.Errstr = "成功写入"
				outstr.Rfid=[]string{"",""}
				b, _ = json.Marshal(outstr)
				fmt.Println("success ",com,outstr)
				w.Write([]byte(b))
				
			}
			
		}

	}
	
}


func jidscan(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
		 
	readid := r.FormValue("readid")
	com := r.FormValue("com")
	var outstr OUTSTRC
	outstr.Errorint = 1
	outstr.Errstr = "串口打不开"
	b, _ := json.Marshal(outstr)
	fmt.Println(com," ",readid)
	c := &serial.Config{Name: com, Baud: 38400}
	s, err := serial.OpenPort(c)
	if err != nil {
		fmt.Println("error ",com)
		w.Write([]byte(b))
		return
	}
	fmt.Println("scan tag id")
	senddata = []byte("\x00\x01\x08\x00\x00\x01\x00\x00")
	onechars, _ := hex.DecodeString(readid)
	senddata[0] = onechars[0]
	senddata[7] = Xor(string(senddata[0:7]))
	n, err := s.Write(senddata)

	if err != nil {
		outstr.Errorint = 1
		outstr.Errstr = "串口发送错误"
		b, _ = json.Marshal(outstr)
		w.Write([]byte(b))
	    return
	}

	buf := make([]byte, 128)
	recvdata := make([]byte, 128)
	num := 0
	for {
		n, err = s.Read(buf)
		for i := 0; i < n; i++ {
			recvdata[i+num] = buf[i]
		}
		num = num + n
		if compbytes(recvdata, num) == 0 {
			fmt.Println(recvdata)
			for i := 0; i < num; i++ {
				strrecv := string(recvdata[i+1 : i+7])
				if strrecv == "000000" {
					break
				}
				fmt.Println(strrecv)
				setonly(strrecv)
				i = i + 6
			}
			break
		}
		if err != nil {
			break
		}
	}
	fmt.Println(ids)
	outstr.Errorint = 0
	outstr.Errstr="成功发送扫描命令"
	outstr.Rfid = ids
	b, _ = json.Marshal(outstr)
	s.Close() 
	
	w.Write([]byte(b))
	if err != nil {
		w.Write([]byte("{error:1 , errstr:'串口读取错误'}"))
        return
	}


	
	 
	
	
}
func main() {
	flag.Parse()

	if packdata() == 1 {
		fmt.Println("pls input parameters")
		flag.Usage()
		return
	}
	
	if srwflag == 3 {
		fmt.Println("请用浏览器打开http://127.0.0.1:8080/src/")
		
		http.HandleFunc("/readrf", jreadrf)
		http.HandleFunc("/writerf", jwriterf)
		http.HandleFunc("/idscan", jidscan)
		http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./html/"))))
		for {
			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				//log.Fatal("ListenAndServer: ", err)
				fmt.Println("ListenAndServer: ", err)

			}
		}
	
    }
	
	
	c := &serial.Config{Name: "COM3", Baud: 38400}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	
	if srwflag == 2 {
		filecontent, err := ioutil.ReadFile(*Input_file)
		filelength := len(filecontent)
		if err != nil {
			log.Fatal(err)
		}
		filecontent1024 := make([]byte, 1024)
		copy(filecontent1024[0:4], []byte(fmt.Sprintf("%04d", filelength)))
		for aa := 0; aa < len(filecontent); aa++ {
			filecontent1024[aa+4] = filecontent[aa]
		}
		filecontent1024[4+filelength] = Xor(string(filecontent1024[:4+filelength]))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(filelength)
		for slicecount := 0; slicecount < 64; slicecount++ {

			senddata[7] = byte(slicecount)

			for ii := 0; ii < 16; ii++ {
				senddata[ii+8] = filecontent1024[slicecount*16+ii]
			}
			senddata[24] = Xor(string(senddata[0:24]))
			fmt.Println(senddata)
		sendagain:

			n, err := s.Write(senddata)

			if err != nil {
				log.Fatal(err)
			}

			buf := make([]byte, 128)
			recvdata := make([]byte, 128)
			num := 0
			for {
				n, err = s.Read(buf)
				for i := 0; i < n; i++ {
					recvdata[i+num] = buf[i]
				}
				num = num + n
				if num >= 2 {
					intnum := int(recvdata[1])
					if num == intnum {
						break
					}
				}

			}
			if recvdata[2] == '\x01' {
				fmt.Println("...")
				time.Sleep(time.Second * 1)
				goto sendagain
			}
			if recvdata[2] == '\x00' {
				fmt.Println("ok", recvdata[7])
			}

		}

	}
	if srwflag == 1 {
		filecontent := make([]byte, 16*64)
		slicecount := 0
		for i0 := 0; i0 < 64; i0++ {
			senddata[7] = byte(i0)
			fmt.Println(senddata[7])
			senddata[8] = Xor(string(senddata[0:8]))
			n, err := s.Write(senddata)

			if err != nil {
				log.Fatal(err)
			}

			buf := make([]byte, 128)
			recvdata := make([]byte, 128)
			num := 0

			for {
				n, err = s.Read(buf)

				for i1 := 0; i1 < n; i1++ {
					recvdata[i1+num] = buf[i1]
				}
				num = num + n
				if num >= 2 {
					intnum := int(recvdata[1])
					if num == intnum {
						fmt.Println(recvdata)
						if readcomppack(recvdata) == 0 || recvdata[num-1] != Xor(string(recvdata[:num-1])) {
							i0 = i0 - 1
							fmt.Println(recvdata[num-1])
							fmt.Println(Xor(string(recvdata[:num-1])))
						} else {
							for i2 := 0; i2 < 16; i2++ {
								filecontent[i2+16*slicecount] = recvdata[i2+4]
							}
							slicecount++

						}

						break
					}
				}

			}

		}
		writefilelength, _ := strconv.Atoi(string(filecontent[0:4]))
		//fmt.Println(writefilelength)
		if filecontent[4+writefilelength] == Xor(string(filecontent[0:4+writefilelength])) {
			ioutil.WriteFile("outputFile"+*Input_read, filecontent[4:writefilelength+4], 0x644)
		} else {
			fmt.Println("file checksum error")
		}
	}
	if srwflag == 0 {
		n, err := s.Write(senddata)

		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 128)
		recvdata := make([]byte, 128)
		num := 0
		for {
			n, err = s.Read(buf)
			for i := 0; i < n; i++ {
				recvdata[i+num] = buf[i]
			}
			num = num + n
			if compbytes(recvdata, num) == 0 {
				fmt.Println(recvdata)
				for i := 0; i < num; i++ {
					strrecv := string(recvdata[i+1 : i+7])
					if strrecv == "000000" {
						break
					}
					fmt.Println(strrecv)
					setonly(strrecv)
					i = i + 6
				}
				break
			}
			if err != nil {
				break
			}
		}
		fmt.Println(ids)
		if err != nil {
			log.Fatal(err)
		}
	}
}

var ids []string
var idsnum = 0

func setonly(str string) {
	for i := 0; i < idsnum; i++ {
		if ids[i] == str {
			return
		}
	}
	ids = append(ids, str)
	idsnum = idsnum + 1

}

func packdata() int {
	termid := *Input_scan
	readid := *Input_read
	writeid := *Input_write
	 
    if (*Input_daemon)[0:3] == "ddd" {
		//fmt.Println((*Input_daemon)[0:3] )
		srwflag = 3
		return 0
	}
	if termid[:len(termid)] != "FFF" {
		if readid != "020304F" || writeid != "020304F" {
			if readid != "020304F" {
				fmt.Println("read tag data")
				senddata = []byte("\x00\x02\x09\x02\x03\x04\x00\x01\x00")
				onechars, _ := hex.DecodeString(termid)
				senddata[0] = onechars[0]
				onechars, _ = hex.DecodeString(readid)
				senddata[3] = onechars[0]
				senddata[4] = onechars[1]
				senddata[5] = onechars[2]
				senddata[8] = Xor(string(senddata[0:8]))
				srwflag = 1
				return 0

			} else {
				fmt.Println("write tag data")
				senddata = make([]byte, 25)
				onechars, _ := hex.DecodeString(termid)
				senddata[0] = onechars[0]
				senddata[1] = '\x03'
				senddata[2] = '\x19'
				onechars, _ = hex.DecodeString(writeid)
				senddata[3] = onechars[0]
				senddata[4] = onechars[1]
				senddata[5] = onechars[2]

				srwflag = 2
				return 0
			}
		} else {
			fmt.Println("scan tag id")
			senddata = []byte("\x00\x01\x08\x00\x00\x01\x00\x00")
			onechars, _ := hex.DecodeString(termid)

			senddata[0] = onechars[0]
			senddata[7] = Xor(string(senddata[0:7]))
			srwflag = 0
			return 0
		}

	}

	return 1
}
func compbytes(bb []byte, num int) int {
	if bb[num-1] != '0' {
		return 1
	}
	if bb[num-2] != '0' {
		return 2
	}
	if bb[num-3] != '0' {
		return 3
	}
	if bb[num-4] != '0' {
		return 4
	}
	if bb[num-5] != '0' {
		return 5
	}
	if bb[num-6] != '0' {
		return 6
	}

	return 0
}
func readcomppack(packbytes []byte) int {
	if packbytes[0] == '\x02' && packbytes[1] == '\x04' && packbytes[2] == '\x01' && packbytes[3] == '\x07' {
		return 0
	}
	return 1
}
func Xor(str string) byte {
	var ret byte
	ret = str[0]
	for i := 0; i < len(str)-1; i++ {
		ret = ret ^ str[i+1]

	}

	return ret

}
