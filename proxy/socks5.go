package main

import(
    "fmt"
    "net"
    "strconv"
    "flag"
)

// socket5协议参考 https://blog.csdn.net/qq_36963214/article/details/115258597
var User = ""
var Passwd = ""

func exchange(a net.Conn,b net.Conn){
    defer a.Close()
    for {
        buff := make([]byte,1024)
        fmt.Println("read")
        n, err := a.Read(buff)
        fmt.Println(n)
        if err != nil {
            fmt.Println(err)
            break
        }
        if n == 0 {
            continue
        }
        fmt.Println("write")
        n1, err1 := b.Write(buff[0:n])
        fmt.Println(n1)
        if err1 != nil {
            break
        }
        if n1 != n {
            fmt.Println(n)
            fmt.Println(n1)
            break
        }
    }
}

func handle_conn(conn net.Conn,check bool){
    defer conn.Close()
    b := make([]byte,1024)

    n, err := conn.Read(b)
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

    go exchange(conn,remote_conn)
    exchange(remote_conn,conn)
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
