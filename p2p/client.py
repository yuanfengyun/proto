import socket
import _thread
import json

severAddr=""
severPort=6000

udpSocket = socket.socket(socket.AF_INET,socket.SOCK_DGRAM)
udpSocket.bind((ip, 30001))#绑定端口
otherAddr=[]#保存其他用户的地址信息

def acc(sc):
    global otherAddr
    cont,destinf=sc.recvfrom(1024)
    jsn=json.loads(cont.decode("utf-8"))
    if(jsn["cod"]==2):#0打洞1初始连接2服务器向其他用户广播新加入地址3向新加入用户通知已在线人员4用户消息
        print(str(jsn['msg'])+"进入聊天")
        otherAddr.append(jsn['msg'])#把新加入用户添加到otherAddr
        sc.sendto('{"cod":0}'.encode('utf-8'),tuple(jsn['msg']))#打洞
    if(jsn["cod"]==3):
        otherAddr=jsn['msg']
        print("进入聊天，其他人有"+str(jsn['msg']))
        for x in jsn['msg']:
            sc.sendto('{"cod":0}'.encode('utf-8'),tuple(x))#打洞
    if(jsn["cod"]==4):
        print(str(destinf)+":"+jsn['msg'])
    acc(sc)
_thread.start_new_thread(acc,(udpSocket,))

jsn={"cod":1,"msg":"a"}
udpSocket.sendto(json.dumps(jsn).encode("utf-8"),(severAddr,severPort))#向服务器第一次请求
 
while(1==1):
    s=input("")
    if(s!=""):
        jsn={"cod":4,"msg":s}
        for x in otherAddr:#向所有其他用户广播
            udpSocket.sendto(json.dumps(jsn).encode("utf-8"),(x[0],x[1]))
