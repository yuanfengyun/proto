package main

import(
    "fmt"
    "net"
    "strconv"
    "flag"
    "io"
)

func handle_conn(conn net.Conn,dest string){
    defer conn.Close()
    remote_conn, err := net.Dial("tcp",dest)
    if err != nil {
        return
    }
    go io.Copy(conn,remote_conn)
    io.Copy(remote_conn,conn)
}

func main(){
    port := flag.Int("l",5000,"listen port")
    dest_port := flag.String("d","6000","dest port")
    dest_host := flag.String("h","127.0.0.1","dest host")

    flag.Parse()

    addr, err := net.ResolveTCPAddr("tcp4", ":" + strconv.FormatInt(int64(*port),10))
    if err != nil {
        return
    }
    fmt.Println(addr)
    l, err1 := net.ListenTCP("tcp", addr)
    if err1 != nil {
        return
    }
    
    dest := *dest_host + ":" + *dest_port
    for {
        conn, err := l.Accept()
        if err != nil {
            continue
        }
        fmt.Println("conn from")
        go handle_conn(conn,dest)
    }
}
