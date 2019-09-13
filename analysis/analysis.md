
## 一次DNS请求分析
Frame 1:  
Internet Protocol Version 4, Src: 127.0.0.1, Dst: 127.0.1.1 //在本机解析域名  
User Datagram Protocol, Src Port: 40723, Dst Port: 53     //UDP协议  
    Source Port: 40723                                    //源端口40723  
    Destination Port: 53                               //目标端口53  
    Length: 41                                         
    Checksum: 0xff3c [unverified]  
    [Checksum Status: Unverified]  
    [Stream index: 0]  
    [Timestamps]   
        [Time since first frame: 0.000000000 seconds]  
        [Time since previous frame: 0.000000000 seconds]   
Domain Name System (query)  
    Transaction ID: 0x3eae         //标识字段,辨别DNS应答报文是哪个请求报文的响应  
    Flags: 0x0100 Standard query   //标志字段  
        0... .... .... .... = Response: Message is a query//1为响应,0为查询      
        .000 0... .... .... = Opcode: Standard query (0)      //0标准,1反向,2服务器状态请求  
        .... ..0. .... .... = Truncated: Message is not truncated    //截断,1：超过512字节截断  
        .... ...1 .... .... = Recursion desired: Do query recursively 
        // 1：得到递归响应  
        .... .... .0.. .... = Z: reserved (0)            //全0保留字段  
        .... .... ...0 .... = Non-authenticated data: Unacceptable//授权回答  
    Questions: 1  
    Answer RRs: 0                   //资源记录数  
    Authority RRs: 0                 //授权资源记录数  
    Additional RRs: 0                //额外资源记录数  
    Queries                       //查询或者响应的正文部分  
        mirror.azure.cn: type A, class IN  
            Name: mirror.azure.cn        //要到达的域名  
            [Name Length: 15]  
            [Label Count: 3]  
            Type: A (Host Address) (1)      //IPv4地址  
            Class: IN (0x0001)  
[Response In: 7]  
   
Frame 2:  
Internet Protocol Version 4, Src: 10.11.29.131, Dst: 210.22.70.3   //路由器向上寻找到达210.22.70.3（本地域名服务器迭代根域名服务器）  
User Datagram Protocol, Src Port: 58768, Dst Port: 53  //源端口58768，目的端口53  
Frame 3:  
Internet Protocol Version 4, Src: 10.11.29.131, Dst:   210.22.84.3//路由器向上寻找到达210.22.84.3（本地域名服务器迭代根域名服务器）  
  
Frame 4:  
Internet Protocol Version 4, Src: 10.11.29.131, Dst:   114.114.114.114//路由器向上寻找到达114.114.114.114（本地域名服务器迭代根域名服务器）  
Frame 5:  
Internet Protocol Version 4, Src: 210.22.84.3, Dst:   10.11.29.131//（跟域名服务器查询到dns向上返回）  
    Questions: 1  
    Answer RRs: 3              // //资源记录数3，说明有三种解析结果  
    Authority RRs: 0  
    Additional RRs: 0  
    Answers                         //dns应答  
        eastmirror.chinacloudapp.cn: type A, class IN, addr 139.217.146.62   //这就是mirror.azure.cn域名对应的地址  
            Name: eastmirror.chinacloudapp.cn  
            Type: A (Host Address) (1)  
            Class: IN (0x0001)  
            Time to live: 60  
            Data length: 4  
            Address: 139.21.146.6  

Frame 6:  
	Internet Protocol Version 4, Src: 210.22.70.3, Dst:  10.11.29.131//另一条解析路径  
            Address: 139.217.146.62    

Frame 7：   
	Internet Protocol Version 4, Src: 127.0.1.1, Dst: 127.0.0.1   //通过递归返回到主机  
eastmirror.chinacloudap.cn: type A, class IN, addr 139.217.146.62  
  
1.	主机向本地域名服务器进行递归查询；  
2.	本地域名服务器以dns客户机身份向三个根域名服务器（210.22.70.3、210.22.84.3、114.114.114.114）发出请求报文  
3.	210.22.84.3和210.22.70.3都查到dns对应的ip，然后向本地服务器返回  
4.	本地服务器把ip返回给主机   


## TCP 握手流程:  
Frame 8:  
No.     Time           Source                Destination           Protocol   Length Info
      8 0.004007447    10.11.29.131          139.217.146.62        TCP      76     60686 → 80 [SYN] Seq=0 Win=29200 Len=0 MSS=1460 SACK_PERM=1 TSval=1759700185 TSecr=0 WS=128

Transmission Control Protocol, Src Port: 60686, Dst Port: 80, Seq: 0, Len: 0  
    Source Port: 60686  
    Destination Port: 80  
    [Stream index: 0]  
    [TCP Segment Len: 0]          //tcp分段序号为0  
    Sequence number: 0    (relative sequence number)  

客户端向服务端发出一次SYN请求，随机序列号seq=x=0（这个序列号是标识报文的开始位置），长度len=0,窗口大小win=29200字节，窗口缩放因子ws=128，实际能接收的大小rwnd=29200*128字节

Frame 9：  
No.     Time           Source                Destination           Protocol Length Info  
      9 0.010020167    139.217.146.62        10.11.29.131          TCP      76     80 → 60686 [SYN, ACK] Seq=0 Ack=1 Win=28960 Len=0 MSS=1440  SACK_PERM=1 TSval=3890471464 TSecr=1759700185 WS=128  

服务端响应客户端请求SYN，ACK ，服务端随机序列号seq=y=0，ack=x+1=1  

Frame 10：  
No.     Time           Source                Destination           Protocol Length Info
     10 0.010069075    10.11.29.131          139.217.146.62        TCP      68     60686 → 80 [ACK] Seq=1 Ack=1 Win=29312 Len=0 TSval=1759700191 TSecr=3890471464  

客户端收到服务端的响应后发送确认信息ACK，seq=x+1+分段序长度0=1，ack=y+1+分段长度0=1，此时完成握手建立连接；

