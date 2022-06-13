package global

var ServerUrl = "http://127.0.0.1:8080/" // 服务器的url

var SqlUsername = "root"                 // 数据库的用户名
var SqlPassword = "123456"               // 数据库的密码
var SqlDBName = "simple_tiktok"          // 使用的数据库名
var SqlDemoDBName = "demo_simple_tiktok" // 导入demo数据使用的数据库名

var SavaFilePath = "./public/"      // 投稿文件保存路径
var MaxFeedVideosNumOnce = 30       // 一次拉取视频最多返回的视频数
var MaxTokenValidTime int64 = 86400 // token的最大有效时间秒数，86400代表一天
var MaxFollowUserCount int64 = -1   // 单个用户最多关注的用户数，-1为无限制
var MaxCommentCount int64 = -1      // 单个视频评论的最大数，-1为无限制
