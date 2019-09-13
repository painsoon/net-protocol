package main

import (
	"fmt"
	"net"
	"strings"
)

const form = `<html><body><form action="#" method="post">
				<input type="text" name="name"/>
				<input type="text" name="psw"/>
				<input type="submit" value="Submit" />
			</form></body></html>`

type Person struct {
	Name  string
	Psw   string
}

type HttpData struct {
	Header string
	body   string
	method string
}

var userMap = make(map[string]Person)

func (httpData *HttpData)DataPara() map[string]string {
	dataMap := make(map[string]string)
	str := httpData.body
	paramnums := strings.Split(str,"&")
	for _,params := range paramnums {
		param := strings.Split(params,"=")
		dataMap[param[0]] = param[1]
	}
	return dataMap
}

func getHttpData(strs []string) *HttpData{
	httpData := &HttpData{}
	var flag = false
	for i,str := range strs{
		if str == "" {
			flag = true
			continue
		}else{
			if i==0 {
				index := strings.Index(str,"/")
				httpData.method = str[:index-1]
			}
		}
		if flag {
			httpData.body = str
		}
	}
	return httpData
}

func HandleConn(conn net.Conn){

	buf := make([]byte,2048)
	n,err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	strs := strings.Split(string(buf[:n]),"\r\n")
	httpData := getHttpData(strs)
	httpData.Header = "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n"

	conn.Write([]byte(httpData.Header))
	switch httpData.method {
	case "GET":
		fmt.Println(len(userMap))
		if len(userMap)>0 {
			conn.Write([]byte("name:"+userMap["user1"].Name))
			conn.Write([]byte("\r\n"))
			conn.Write([]byte("password:"+userMap["user1"].Psw))
		}else {
			conn.Write([]byte(form))
		}

	case "POST":
		fmt.Println(httpData.body)
		dataMap := httpData.DataPara()
		person := Person{dataMap["name"],dataMap["psw"]}
		userMap["user1"] = person
		conn.Write([]byte("name:"+person.Name))
		conn.Write([]byte("\r\n"))
		conn.Write([]byte("password:"+person.Psw))
	}

	defer conn.Close()

}

func main(){
	listen,err := net.Listen("tcp","127.0.0.1:8080")
	defer listen.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	for{
		conn,err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go HandleConn(conn)
	}

}

