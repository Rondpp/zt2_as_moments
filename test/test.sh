# 获取个人信息
curl  -H "token:token10001" -H "accid:10001" "http://47.98.246.74:8002/user/?accid=10001"

#修改个人信息
curl  -H "token:token10001" -H "accid:10001" -d '{"name":"英雄10001","avatar":"https://zt2as.oss-cn-hangzhou.aliyuncs.com/image/2018-05-23/240603656/1527087249223-82FAA794-F50F-41B3-8BE3-5A92E0DEFADD.jpeg", "sex":1,"birthday":1212141241}'  "http://47.98.246.74:8002/user/?accid=10001"

#发布个人动态
curl  -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态","pic":["https://zt2as.oss-cn-hangzhou.aliyuncs.com/image/2018-05-23/240603656/1527087249223-82FAA794-F50F-41B3-8BE3-5A92E0DEFADD.jpeg","https://zt2as.oss-cn-hangzhou.aliyuncs.com/image/2018-05-23/240603656/1527087249227-D0B8F08E-5745-48C9-9662-5EA16A4FFEA9.jpeg"]}' http://47.98.246.74:8002/moments/

#发布视频动态
curl  -H "token:token10002" -H "accid:10002"  -d '{"content":"我的第一条动态","video":"68c233deb5c64071b78d7acab4f61169"}' http://47.98.246.74:8000/moments/

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

#获取动态的评论 
curl   -H "token:token10002" -H "accid:10002"   "http://47.98.246.74:8002/comment/?moment_id=5aeef754e71e163e45d6e83a&start_id=&num=10"

#获取评论的评论
curl   -H "token:token10002" -H "accid:10002"   "http://47.98.246.74:8002/comment/?comment_id=5aeef754e71e163e45d6e83a&start_id=&num=10"

#获取单个评论
curl   -H "token:token10002" -H "accid:10002"   "http://47.98.246.74:8002/comment/?comment_id=5aeef754e71e163e45d6e83a"


#获取我相关的评论
curl   -H "token:token10001" -H "accid:10001"   "http://47.98.246.74:8002/message/?type=1&start_id=&num=10"

#删除我相关的评论
curl   -H "token:token10001" -H "accid:10001"  -X DELETE  "http://47.98.246.74:8002/message/?type=1&id=5aeef754e71e163e45d6e83a"

#查询是否有未读消息
curl   -H "token:token10001" -H "accid:10001"  "http://47.98.246.74:8002/message/unread?"

#设置权限
curl   -H "token:token10000" -H "accid:10000"  -d '{"accid":10002, "account":"", "permission":15}' "http://47.98.246.74:8000/admin/user/?"

#获取管理员列表
curl   -H "token:token10000" -H "accid:10000"  "http://47.98.246.74:8000/admin/user/?"

#删除权限
curl   -H "token:token10000" -H "accid:10000"  -X DELETE  "http://47.98.246.74:8000/admin/user/?accid=10002"

#禁言
curl   -H "token:token10002" -H "accid:10002"  -d '{"accid":10003, "time":1000}' "http://47.98.246.74:8000/admin/forbidden/?"

#置顶
curl   -H "token:token10002" -H "accid:10002"  -d '' "http://47.98.246.74:8000/admin/totop/?moment_id=5b03996ae71e16447a5bacc5"

#删除动态
curl   -H "token:token10002" -H "accid:10002"  -d '' "http://47.98.246.74:8000/admin/delete/?moment_id=5b03996ae71e16447a5bacc5"

#删除评论
curl   -H "token:token10002" -H "accid:10002"  -d '' "http://47.98.246.74:8000/admin/delete/?comment_id=5afd3d12e71e16741b960397"


#获取需要审核的视频列表
curl  -H "token:token10002" -H "accid:10002"    "http://47.98.246.74:8000/admin/moments/?start_id=&num=10"

#审核视频
curl  -H "token:token10002" -H "accid:10002"  -d ''   "http://47.98.246.74:8000/admin/moments/?moment_id=5b0bb04be71e166d647636f7&pass=1"
