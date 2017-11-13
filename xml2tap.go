package main

import (
	"log"
	"gopkg.in/natefinch/lumberjack.v2"//rotational logging
	"net"
	"tap"
	//"strings"
	"flag"
	"encoding/xml"
	//"os"
	"net/http"
	"fmt"
	//"html/template"
	"time"
)
type Page struct {
	ID string `xml:"ID"`
	TagText string `xml:"TagText"`
	//Type string `xml:"Type"`
}

var queuesize = 0//the size of the processed message channel

func HomePage(w http.ResponseWriter, req *http.Request) {
	
}

func StatusPage(w http.ResponseWriter, req *http.Request) {
	if queuesize <=300{
		fmt.Fprintf(w,"OK: Curent Queue Size:%v",queuesize)
	}else{
		fmt.Fprintf(w,"ERROR: Curent Queue Size:%v",queuesize)
	}
	
}

func SendPage(w http.ResponseWriter, req *http.Request) {

}

func webserver(msgchan chan string,portnum string) {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/status", StatusPage)
	http.HandleFunc("/page", SendPage)
	for{
		log.Print(http.ListenAndServe(":"+portnum, nil))
	}
}

func queuemonitor(msgchan chan string){
	for{
	queuesize = len(msgchan)
	time.Sleep(5 * time.Second)
	}
}

//example r5 xml
//<Page xmlns:xsi='http://www.w3.org/2001/XMLSchema-instance' xmlns:xsd='http://www.w3.org/2001/XMLSchema'>
//	<ID>89699</ID>
//	<TagText>4906 beeping</TagText>
//   <Type>Phone/Pager</Type>
//</Page>

//example response for a ___PING___
//<?xml version="1.0" encoding="utf-8"?> <PageTXSrvResp State="7" PagesInQueue="4" PageOK="1" />

//main can accept 2 flag arguments the port for the xml listener and the port
//for the TAP output listener
//i.e call xml2tap -xmlPort=5051 tapPort=10001
//default ports are 5051 for xml and 10001 for tap
func main() {
	

	xmlPort := flag.String("xmlPort","5051","xml listener port for localhost")
	tapPort := flag.String("tapPort","10001","localhost listener port for TAP server")
	httpPort := flag.String("httpPort","80","localhost listner port for http server")
	flag.Parse()
	
	// //log output file
	// f, err := os.OpenFile("filename", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// if err != nil {
	// 		log.Fatal(err)
	// }   
	// //defer to close 
	// defer f.Close()

	// //set output of logs to file
	// log.SetOutput(f)


	log.SetOutput(&lumberjack.Logger{
		Filename:   "/var/log/xml2tap/xml2tap.log",
		MaxSize:    100, // megabytes
		MaxBackups: 5,
		MaxAge:     60, //days
		Compress:   true, // disabled by default
	})

	log.Printf("STARTING XML Listener on tcp port %v\n\n", *xmlPort)
	l, err := net.Listen("tcp", ":"+ *xmlPort)
	if err != nil {
		log.Println("Error opening XML listener, check log for details")
		log.Fatal(err)
	}
	defer l.Close()
		
	//message processing channel for sip2tap conversions 1000 chosen as 1 per bed maximum
 	parsedmsgs := make(chan string, 1000)

	//start a tap server
	go tap.Server(parsedmsgs,*tapPort)
	
	//start a webserver
	go webserver(parsedmsgs,*httpPort)
	
	//start a queue monitor
	go queuemonitor(parsedmsgs)
	

	for {
		
		// Listen for an incoming xml connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			log.Fatal(err)
		}

		// Handle connections in a new goroutine.
		go func(c net.Conn, msgs chan<- string) {
			//set up a decoder on the stream
			d := xml.NewDecoder(c)
			
			for {
				// Look for the next token
				// Note that this only reads until the next positively identified
				// XML token in the stream
				t, err := d.Token()
				if err != nil  {
					log.Printf("Token error %v",err.Error())
					break
				}
				switch et := t.(type) {
				
					case xml.StartElement:
						// search for Page start element and decode
						if et.Name.Local == "Page" {
							p := &Page{}
							// decode the page element while automagically advancing the stream
							// if no matching token is found, there will be an error
							// note the search only happens within the parent.
							if err := d.DecodeElement(&p, &et); err != nil {
								log.Printf("error decoding element %v",err.Error())
								c.Close()
								return
							}
				
							// We have decoded the xml message now send it off to TAP server or reply if ping
							log.Printf("Pin:%v;Msg:%v\n", p.ID,p.TagText)

							//note the R5 system periodically sends out a PING looking for a response
							//this will handle that response or put the decoded xml into the TAP output queue
							if p.ID == "" && p.TagText =="___PING___"{
								//send response to connection
								response := "<?xml version=\"1.0\" encoding=\"utf-8\"?> <PageTXSrvResp State=\"7\" PagesInQueue=\"0\" PageOK=\"1\" />"
								log.Printf("Responding:%v\n",response)
								c.Write([]byte(response))
							}else{
								parsedmsgs <-string(p.ID) + ";" + string(p.TagText)
								
							}


						} 
				
					case xml.EndElement:
						if et.Name.Local != "Page" {
							continue
						}
				}
				
			}	
				
		c.Close()	
		}(conn, parsedmsgs)
	}

}