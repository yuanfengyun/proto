package main

import(
    "fmt"
    "net"
    "strconv"
    "flag"
    "io"
    "net/url"
    "bytes"
    "strings"
)

// socket5协议参考 https://blog.csdn.net/qq_36963214/article/details/115258597
var User = ""
var Passwd = ""

func handle_http(conn net.Conn,b []byte,n int){
        defer conn.Close()

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

func handle_conn(conn net.Conn,check bool){
    b := make([]byte,1024)

    n, err := conn.Read(b)
    if err != nil || n == 0 {
        return
    }
    if int(b[0]) != 5 {
        handle_http(conn,b,n)
        return
    }

    defer conn.Close()
    if err != nil || n < 3 || n > 257 {
        return
    }

    ver, nmethods := int(b[0]),int(b[1])
    if ver != 0x05 || nmethods + 2 != n {
        return
    }

    if check {
        conn.Write([]byte("\x05\x01"))
        n, err = conn.Read(b)

        if err != nil {
            fmt.Println(b)
            return
        }

        ulen := int(b[1])
        if n < ulen + 4 {
            fmt.Println("4")
            return
        }

        uname := string(b[2:2+ulen])
        plen := int(b[2 + ulen])
        if plen + ulen + 3 != n{
            fmt.Println("5")
            return
        }

        passwd := string(b[2 + ulen + 1: 2+ulen+1+plen])

        if User != uname || passwd != Passwd {
            conn.Write([]byte("\x05\x01"))
            fmt.Println("6")
            return
        }

        conn.Write([]byte("\x05\x00"))


    }else{
        conn.Write([]byte("\x05\x00"))
    }
    n, err = conn.Read(b)

    if err != nil || n < 10 {
        fmt.Println("7")
        return
    }

    cmd := int(b[1])
    if cmd != 1 {
        fmt.Println("8")
        return
    }
    atyp := int(b[3])

    var remote_conn net.Conn
    var remote_addr string
    var port int
    var end int
    // ipv4
    if atyp == 1 {
        remote_addr = string(b[4:8])
        port = int(b[9]) << 8 + int(b[10])
        end = 11
    }else if atyp == 3 {
        hostlen := int(b[4])
        if n != hostlen + 7{
            return
        }
        remote_addr = string(b[5:5+hostlen])
        port = (int(b[5+hostlen]) << 8) + int(b[5+hostlen+1])
        end = 5+hostlen+2
    }else{
        fmt.Println("10")
        return
    }
    remote_conn, err1 := net.Dial("tcp", remote_addr + ":" + strconv.FormatInt(int64(port),10))
    if err1 != nil {
        fmt.Println("11")
        return
    }
    b[1] = 0
    conn.Write(b[:n])
    fmt.Println(remote_addr)
    if n > end {
        remote_conn.Write(b[end+1:])
        fmt.Println("haha")
    }

    go io.Copy(conn,remote_conn)
    io.Copy(remote_conn,conn)
}

func main(){
    port := flag.Int("p",5000,"port")
    user := flag.String("user","user","user")
    passwd := flag.String("passwd","passwd","password")

    flag.Parse()

    User = *user
    Passwd = *passwd

    addr, err := net.ResolveTCPAddr("tcp4", ":" + strconv.FormatInt(int64(*port),10))
    if err != nil {
        return
    }
    fmt.Println(addr)
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
        go handle_conn(conn,false)
    }
}
