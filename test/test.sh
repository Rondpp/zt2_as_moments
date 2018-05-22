# 获取个人信息
curl  -H "token:token10001" -H "accid:10001" "http://47.98.246.74:8002/user/?accid=10001"

#修改个人信息
curl  -H "token:token10001" -H "accid:10001" -d '{"name":"英雄10001","avatar":"http://xxxx.xxx/", "sex":1,"birthday":1212141241}'  "http://47.98.246.74:8002/user/?accid=10001"

#发布个人动态
curl  -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态","pic":["http://47.98.246.74:8001/pic/2.jpg","http://47.98.246.74:8001/pic/3.jpg"]}' http://47.98.246.74:8002/moments/

#发布视频动态
curl  -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态","video":"https://zt2as.oss-cn-hangzhou.aliyuncs.com/abcd/qwe/vide-o4.mp4"}' http://47.98.246.74:8000/moments/

#删除个人动态
curl  -H "token:token10001" -H "accid:10001"  -X DELETE http://47.98.246.74:8002/moments/?moment_id=5aed8236e71e16138e1f57d8

#获取个人动态
curl  -H "token:token10001" -H "accid:10001"  "http://47.98.246.74:8002/moments/?accid=10002&type=0&start_id=&num=10"

#获取最新动态
curl  -H "token:token10001" -H "accid:10001"  "http://47.98.246.74:8002/moments/?type=1&start_id=&num=10"

#获取最新视频
curl  -H "token:token10001" -H "accid:10001"  "http://47.98.246.74:8002/moments/?type=2&start_id=&num=10"

#获取关注的人的动态和视频
curl  -H "token:token10001" -H "accid:10001"  "http://47.98.246.74:8002/moments/?type=3&start_id=&num=10"

#获取某个动态详情
curl  -H "token:token10001" -H "accid:10001"  "http://47.98.246.74:8000/moments/?type=4&moment_id=5afce70de71e163c8888044c"

#关注
curl  -H "token:token10002" -H "accid:10002" -d""  http://47.98.246.74:8002/follow/?accid=10001

#取关
curl  -H "token:token10002" -H "accid:10002" -X DELETE  http://47.98.246.74:8002/follow/?accid=10001

#获取关注的人
curl  -H "token:token10002" -H "accid:10002"   "http://47.98.246.74:8002/follow/?accid=10002&start_id=&num=10"

#获取粉丝
curl  -H "token:token10002" -H "accid:10002"   "http://47.98.246.74:8002/fans/?accid=10002&start_id=&num=10"

#搜索动态
curl   -H "token:token10002" -H "accid:10002"  "http://47.98.246.74:8002/search/?keyword=我&type=1&start_id=&num=10"

#搜索人
curl   -H "token:token10002" -H "accid:10002"  "http://47.98.246.74:8002/search/?keyword=100&type=2&start_id=&num=10"

#点赞
curl   -H "token:token10002" -H "accid:10002" -d ""  "http://47.98.246.74:8002/like/?moment_id=5aed8236e71e16138e1f57d8"

#评论 
curl   -H "token:token10002" -H "accid:10002" -d '{"moment_id":"5aeef754e71e163e45d6e83a","ref_comment_id":2,"content":"我的评论2"}'  "http://47.98.246.74:8002/comment/?"

#获取评论 
curl   -H "token:token10002" -H "accid:10002"   "http://47.98.246.74:8002/comment/?moment_id=5aeef754e71e163e45d6e83a&start_id=&num=10"

#获取我相关的评论
curl   -H "token:token10001" -H "accid:10001"   "http://47.98.246.74:8002/message/?type=1&start_id=&num=10"

#删除我相关的评论
curl   -H "token:token10001" -H "accid:10001"  -X DELETE  "http://47.98.246.74:8002/message/?type=1&id=5aeef754e71e163e45d6e83a"

#设置权限
curl   -H "token:token10000" -H "accid:10000"  -d '{"accid":10002, "account":"", "permissions":15}' "http://47.98.246.74:8000/admin/user/?"

#获取管理员列表
curl   -H "token:token10000" -H "accid:10000"  "http://47.98.246.74:8000/admin/user/?"

#删除权限
curl   -H "token:token10000" -H "accid:10000"  -X DELETE  "http://47.98.246.74:8000/admin/user/?accid=10002"
