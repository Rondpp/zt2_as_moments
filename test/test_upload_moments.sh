curl -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态10002_1","pic":["http://47.98.246.74:8001/pic/2.jpg","http://47.98.246.74:8001/pic/3.jpg"]}' http://47.98.246.74:8002/moments/
curl -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态10002_1","pic":["http://47.98.246.74:8001/pic/2.jpg","http://47.98.246.74:8001/pic/3.jpg"]}' http://47.98.246.74:8003/moments/

curl -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态10002_2","pic":["http://47.98.246.74:8001/pic/2.jpg","http://47.98.246.74:8001/pic/3.jpg"]}' http://47.98.246.74:8002/moments/
curl -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态10002_2","pic":["http://47.98.246.74:8001/pic/2.jpg","http://47.98.246.74:8001/pic/3.jpg"]}' http://47.98.246.74:8003/moments/

curl -H "token:token10001" -H "accid:10001"  -d '{"content":"我的第一条动态10001_2","pic":["http://47.98.246.74:8001/pic/2.jpg","http://47.98.246.74:8001/pic/3.jpg"]}' http://47.98.246.74:8002/moments/
curl -H "token:token10001" -H "accid:10001"  -d '{"content":"我的第一条动态10001_2","pic":["http://47.98.246.74:8001/pic/2.jpg","http://47.98.246.74:8001/pic/3.jpg"]}' http://47.98.246.74:8003/moments/

curl -H "token:token10002" -H "accid:10002"  -d '{"content":"视频动态_10002","video":"https://zt2as.oss-cn-hangzhou.aliyuncs.com/abcd/qwe/vide-o4.mp4"}' http://47.98.246.74:8002/moments/
curl -H "token:token10002" -H "accid:10002"  -d '{"content":"视频动态_10002","video":"https://zt2as.oss-cn-hangzhou.aliyuncs.com/abcd/qwe/vide-o4.mp4"}' http://47.98.246.74:8003/moments/

