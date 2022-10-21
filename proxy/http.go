package main

import(
        "fmt"
        "net"
        "strconv"
        "flag"
"bytes"
"strings"
"io"
"net/url"
)

// 协议参考 https://blog.csdn.net/weixin_43507410/article/details/124839308

func handle(conn net.Conn){
        defer conn.Close()

        b := make([]byte,1024)
        n,err := conn.Read(b)
        if err != nil {
                return
        }

        var method, targeturl, address string
        fmt.Sscanf(string(b[:bytes.IndexByte(b[:], '\n')]), "%s%s", &method, &targeturl)

        hostPort, err := url.Parse(targeturl)
        if err != nil {
                return
        }

        if method == "CONNECT" {
                address = hostPort.Scheme + ":" + hostPort.Opaque
        } else {
                address = hostPort.Host
                if strings.Index(hostPort.Host, ":") == -1 { //host 不带端口， 默认 80
                        address = hostPort.Host + ":80"
                }
        }

        //获得了请求的 host 和 port，向服务端发起 tcp 连接
        server, err := net.Dial("tcp", address)
        if err != nil {
                return
        }

        //如果使用 https 协议，需先向客户端表示连接建立完毕
        if method == "CONNECT" {
                fmt.Fprint(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
        } else { //如果使用 http 协议，需将从客户端得到的 http 请求转发给服务端
                server.Write(b[:n])
        }

        go io.Copy(server, conn)
        io.Copy(conn, server)
}

func main(){
        port := flag.Int("p",5000,"port")

        flag.Parse()

        addr, err := net.ResolveTCPAddr("tcp4", ":" + strconv.FormatInt(int64(*port),10))
        if err != nil {
                return
        }
        l, err1 := net.ListenTCP("tcp", addr)
        if err1 != nil {
                return
        }
        for {
                conn, err := l.Accept()
                if err != nil {
                        continue
                }
                fmt.Println("conn from")
                go handle(conn)
        }
}
