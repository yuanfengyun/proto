import socket
import _thread
import json
import time
udpSocket = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
udpSocket.bind(("0.0.0.0", 6000))#注意如果使用云服务器，此处应绑定云服务器内网ip
add=[]#保存所用连接用户

def acc(sc):
    cont,destinf=sc.recvfrom(1024)
    print(cont.decode("utf-8"))
    print(destinf)
    jsn=json.loads(cont.decode("utf-8"))
    if(jsn["cod"]==1):#用户连接
        if(len(add)>0):
            for x in add:
                print(x)
                jtt={"cod":2,"msg":destinf}
                sc.sendto(json.dumps(jtt).encode("utf-8"),x)#向所有用户广播新加入用户的地址
            jtn={"cod":3,"msg":add}
            sc.sendto(json.dumps(jtn).encode('utf-8'),destinf)
        add.append(destinf)
    acc(sc)
_thread.start_new_thread(acc,(udpSocket,))
 
 
s=input("")
while(1==1):
    time.sleep(10) 
