package global

var ServerUrl = "http://127.0.0.1:8080/" // 服务器的url

var SqlUsername = "root"                 // 数据库的用户名
var SqlPassword = "123456"               // 数据库的密码
var SqlDBName = "simple_tiktok"          // 使用的数据库名
var SqlDemoDBName = "demo_simple_tiktok" // 导入demo数据使用的数据库名

var MaxFollowUserCount int64 = 100000000 // 单个用户最多关注的用户数
var MaxCommentCount int64 = 100000000000 // 单个视频评论的最大数
