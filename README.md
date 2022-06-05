# simple_tiktok

## 项目基本信息

### 文件夹结构

```
simple_tiktok
|-- controller
  |-- comment.go    // 评论功能
  |-- common.go     // 使用到的数据结构体和全局服务器url
  |-- demo_data.go  // demo数据，包括导入方法
  |-- favorite.go   // 点赞功能
  |-- feed.go       // 拉去视频
  |-- publish.go    // 发布视频和作品列表
  |-- relation.go   // 关注功能
  |-- user.go       // 用户登录注册、获取用户信息
|-- public
  |-- bear.mp4      // 本地视频文件，熊的视频
  |-- data
|-- repository      // 存储相关
  |--db_init.go     // 与数据库交互相关的数据模型和方法
|-- main.go
|-- router.go
```

### 运行方式

#### 运行环境

+ 服务端：Windows\Linux\Mac
+ 客户端：安卓虚拟机
+ Go版本：Go 1.17+
+ 数据库：mysql 8.0+

#### 运行步骤

1. 设置服务器url和数据库相关信息

首先在`controller/common.go`文件中，修改变量`serverUrl`的值，设置服务器的url，如`var serverUrl = "http://127.0.0.1:8080/"`

然后在`repository/db_init.go`文件中，设置连接数据库的用户名`SqlUsername`、密码`SqlPassword`、数据库名`SqlDBName`，如果需要直接导入demo数据的话，设置导入demo数据使用的数据库名`SqlDemoDBName`（使用Demo数据会先清空之前的数据库再导入数据，如果想使用之前程序保存在数据库中的数据，需要使用名`SqlDBName`的普通数据库）

2. 启动服务端

在终端进入`simple_tiktok`文件目录下，输入

```
go run main.go router.go
```

如果需要使用demo数据库并导入demo数据，输入

```
go run main.go router.go --demo
```

即可在默认的8080端口运行起来服务端

3. 启动客户端

在安卓虚拟机中安装`app-release.apk`，然后在客户端打开软件，连续点击右下方的“我”2~3次，弹出“高级设置”界面，设置服务器url，和之前在`controller/common.go`文件中设置的服务器url一样，如`http://127.0.0.1:8080/`，点击保存并重启，退出软件，再点击进入，即可完成相应设置

![](./imgs/1.png)

### 功能说明

+ 登录、注册功能：每次启动程序会将用户信息导入到内存中，方便其它功能查询，注册会同时更新内存和数据库中的用户信息
+ 视频 Feed 流：支持所有用户刷抖音，包括未登录的用户，视频按照投稿时间倒序推出，单词最多30个
+ 视频投稿：要求用户必须处在登录状态，可以选择自己拍的视频上传，视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可
+ 个人信息：登录用户可以在个人信息页查看自己的用户信息，包括用户名、关注数和粉丝数等
+ 点赞功能：登录用户可以对视频点赞或取消点赞，在个人主页能够查看点赞视频列表，未登录状态下也可以看到视频的总点赞数
+ 评论功能：登录用户可以对视频评论，未登录状态下可以看到视频的总评论数，也可以查看所有用户的评论，按发布时间倒序
+ 关注功能：登录用户可以关注或取消关注其他用户，可以在个人信息页点击打开关注列表和粉丝列表

### 实现思路

#### 1. 登录、注册功能

登录注册过程中使用到的与数据库进行交互的用户信息结构体`UserDao`仿照程序使用的`controller/commont.go`文件中的用户结构体`User`进行设计，放在`repository/db_init.go`文件中，设置用户ID为自增主键：

```go
// 用户信息表users
type UserDao struct {
	Id            int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	Token         string `json:"token"`
}
```

其次有两个map，即`usersLoginInfo`和`userIdToToken`，分别存储用户token与用户User结构体的对应关系和存储用户Id与用户token的对应关系，每次启动程序时，都会将数据库中的用户信息导入这两个map中，可以加快其它功能的查询：

注册时，首先判断用户是否存在，如果用户已经存在，直接返回注册失败，否则将注册用户的用户名填入`Name`属性，用户名加密码组成`Token`属性，`FollowCount`、`FollowerCount`为0，`IsFollow`为`false`构成的UserDao对象直接插入数据库中的`users`表中，并在`usersLoginInfo`和`userIdToToken`记录新注册用户的对应关系。

登录时，首先根据用户名和密码组成token，然后在`usersLoginInfo`中搜索是否存在这个token，存在则可以从`usersLoginInfo`中根据token取出用户信息并返回，找不到token则返回用户不存在。

#### 2. 视频 Feed 流

拉取视频流过程中使用到的与数据库进行交互的用户信息结构体`VideoDao`仿照程序使用的`controller/commont.go`文件中的视频结构体`Video`进行设计，`Video`的`User`类型字段`Author`改成了`int64`类型的字段`AuthorId`，在数据库中只存储视频作者的用户Id，当要获取这个作者的具体用户信息时，通过`userIdToToken`首先将用户Id转成用户token，然后再用token从`usersLoginInfo`中找到具体的用户信息，同时新加了一个`PublishTime`字段用来表示投稿时间，放在`repository/db_init.go`文件中，设置视频ID为自增主键：

```go
// 视频信息表videos
type VideoDao struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	AuthorId      int64  `json:"author_id,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
	PublishTime   int64  `json:"publish_time,omitempty"`
}
```

当拉取视频时，首先获取限制返回视频的投稿时间戳`latest_time`，如果没有设置，则默认为当前时间，然后找到投稿时间不晚于lastestTime的投稿视频，按投稿时间倒序排列，最多30个，接着获取用户点赞的视频列表，并存储在map中，根据投稿视频是否在这个点赞的视频列表中，设置获取的投稿视频的是否点赞`IsFavorite`属性，并记录本次返回的视频中，发布最早的时间`nextTime`，作为下次请求时的`latest_time`

#### 3. 视频投稿

当投稿时，首先获取token，判断是否处于登录状态，不是则直接返回，取消发布视频，否则获取发布的视频文件数据，将视频文件经过文件路径和命名处理后存入本地的`public`文件夹目录下，同时将相应的视频信息组成一个`VideoDao`插入数据库的`videos`表中，其中播放路径`PlayUrl`由放在`common.go`文件中的服务器url、`static/`、文件最终名字`finalName`组成，发布时间通过`time.Now().Unix()`方法获取，即为当前时间。

当获取发布作品列表时，首先从数据库中的`videos`表根据用户id获取这个用户发布的视频列表，因为这个用户有可能对自己的视频点赞，所以也要获取这个用户点赞的视频列表，设置视频列表中的视频的`IsFavorite`是否点赞属性后返回这个用户发布视频列表。

#### 4. 个人信息

在登录时会获取个人信息，首先在`usersLoginInfo`中查找token，如果找不到返回用户不存在，找到则从`usersLoginInfo`中根据token取出用户信息，因为每个用户的关注列表不同，所以每次获取个人信息时，都需要从数据库的`follows`表中，找到当前用户关注的所有用户的Id，然后根据是否关注，重新设置`usersLoginInfo`中每个存在用户的`IsFollow`属性，被当前用户关注则`IsFollow`属性为`true`，否则为`false`，更新当前用户的关注状态，同时还要获取这个用户的视频列表，更新这个用户对每个视频的点赞状态。

#### 5. 点赞功能

点赞或取消点赞过程中使用到的与数据库进行交互的用户信息结构体是`FavoriteVideoDao`有两个字段，分别是给这个视频点赞的用户token和这个视频的id，放在`repository/db_init.go`文件中：

```go
// 点赞视频信息表favorite_videos
type FavoriteVideoDao struct {
	Token   string `json:"token"`    // 用户的token
	VideoId int64  `json:"video_id"` // 用户喜欢的视频Id
}
```

当点赞或取消点赞时，首先获取token，判断是否处于登录状态，不是则直接返回，否则先根据token获取这个用户点赞的视频列表，判断用户是否对当前视频点过赞了，如果是点赞行为且之前没有给这个视频点过赞，那么将数据库中`videos`表的这个视频的总点赞数加一，并在`favorite_videos`点赞视频表中创建相应点赞记录，之前点过赞则不作反应；如果是取消点赞行为且之前给这个视频点过赞了，那个更新数据库中的视频总点赞数，删除点赞记录，之前没点过赞则不作反应。

当获取点赞列表时，首先根据用户Id从`userIdToToken`中获取用户token，然后根据用户token从数据库的`favorite_videos`表中获取这个用户点赞的所有视频，最后返回。

#### 6. 评论功能

评论或取消评论过程中使用到的与数据库进行交互的用户信息结构体`CommentDao`仿照程序使用的`controller/commont.go`文件中的评论结构体`Comment`进行设计，`Comment`的`User`类型字段`User`改成了`int64`类型的字段`UserId`，在数据库中只存储评论者的Id，当要获取这个评论者的具体用户信息时，通过`userIdToToken`首先将用户Id转成用户token，然后再用token从`usersLoginInfo`中找到具体的用户信息，同时新加了`VideoId`字段表示视频，`PublishTime`字段用来表示评论时间，放在`repository/db_init.go`文件中，设置评论ID为自增主键：

```go
// 评论信息表comments
type CommentDao struct {
	Id          int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserId      int64  `json:"user_id"`
	VideoId     int64  `json:"video_id"`
	Content     string `json:"content"`
	CreateDate  string `json:"create_date"`
	PublishTime int64  `json:"publish_time"`
}
```

当进行评论或取消评论时，首先获取token，判断是否处于登录状态，不是则直接返回，否则先根据token获取这个用户点赞的视频列表，判断用户是否对当前视频点过赞了，如果是点赞行为且之前没有给这个视频点过赞，那么将数据库中`videos`表的这个视频的总点赞数加一，并在`favorite_videos`点赞视频表中创建相应点赞记录，之前点过赞则不作反应；如果是取消点赞行为且之前给这个视频点过赞了，那个更新数据库中的视频总点赞数，删除点赞记录，之前没点过赞则不作反应。

当获取视频的所有评论时，从数据库的`comments`表中根据视频id获取按发布时间倒序的所有评论，然后从从记录账号信息的`usersLoginInfo`和`userIdToToken`根据用户Id获取具体的评论用户信息，设置获取的评论列表中评论作者字段，返回相应的评论列表。

#### 7. 关注功能

关注或取消关注过程中使用到的与数据库进行交互的用户信息结构体是`FollowDao`有两个字段，分别是给这个关注者的用户id和被关注者的用户id，放在`repository/db_init.go`文件中：

```go
// 关注信息表follows
type FollowDao struct {
	UserId   int64 `json:"user_id"`     // 关注者的用户Id
	ToUserId int64 `json:"to_user_id"`  // 被关注者的用户Id
}
```

当进行关注或取消关注时，首先获取token，判断是否处于登录状态，不是则直接返回，否则先根据token获取这个用户的Id，如果是关注行为，那么在数据库的`follows`表中创建相应的记录，在`users`表中更新关注用户和被关注用户的关注数和被关注数；如果是取消关注行为，那么在数据库的`follows`表中删除相应的记录，在`users`表中更新关注用户和被关注用户的关注数和被关注数，对数据库操作完后，同时也要更新内存中的存储账号信息`usersLoginInfo`的map中相应用户的关注数、被关注数、是否被关注等信息。

当获取关注列表或粉丝列表时，从数据库的`follows`表中根据用户Id获取所有这个用户关注或关注这个用户的用户Id，然后从从记录账号信息的`usersLoginInfo`和`userIdToToken`根据用户Id获取具体的关注或粉丝用户信息，然后返回获取的用户列表。

### 功能展示

这里导入Demo数据进行展示，Demo数据写在`controller/demo_data.go`文件中，使用数据库`demo_simple_tiktok`存放导入的Demo数据，在服务器终端输入`go run main.go router.go --demo`即可导入Demo数据，在Demo数据中：

+ 有5个用户user1、user2、user3、user4、user5，用户名和Token分别为：
  + 用户名：user1@1.com，Token：user1@1.com111111
  + 用户名：user2@2.com，Token：user2@2.com222222
  + 用户名：user3@3.com，Token：user3@3.com333333
  + 用户名：user4@4.com，Token：user4@4.com444444
  + 用户名：user5@5.com，Token：user5@5.com555555
+ 有31个视频，视频Id从1到31，每个视频的发布时间也是从1到31，与Id一样，第一个Id为1的视频是Bear.mp4熊视频，存放在本地，其它视频的播放路径都来自网络，节省仓库存储空间，每个视频的作者分别为：
  + Id为1（熊视频）、6、11、16、21、26的视频作者为user1
  + Id为2（橙子视频）、7、12、17、22、27的视频作者为user2
  + Id为3、8、13、18、23、28的视频作者为user3
  + Id为4、9、14、19、24、29的视频作者为user4
  + Id为5、10、15、20、25、30、31（甜甜圈视频）的视频作者为user5
+ 视频点赞情况为：
  + Id为31的视频（甜甜圈视频）有user3、user4、user5三个用户点赞
  + Id为30的视频有user1、user2两个用户点赞
  + Id为29的视频有user3、user4、user5三个用户点赞
  + Id为2的视频（橙子视频）有user5一个用户点赞
  + Id为1的视频（熊视频）有user1、user2、user4三个用户点赞
+ 视频评论情况为：
  + Id为31的视频（甜甜圈视频）有user1、user2、user3、user4、user5五个用户评论：
    + user1在时间50评论
    + user2在时间49评论
    + user3在时间48评论
    + user4在时间47评论
    + user5在时间46评论
  + Id为30的视频有user1、user2两个用户评论：
    + user3在时间45评论
    + user4在时间44评论
  + Id为2的视频（橙子视频）有user3、user4两个用户评论：
    + user3在时间43评论
    + user4在时间42评论
  + Id为1的视频（熊视频）有user5一个用户评论：
    + user5在时间41评论
+ 关注情况为：
  + user1关注了：user2
  + user2关注了：user3
  + user3关注了：user4
  + user4关注了：user1、user5
  + user5关注了：user1

#### 1. 登录、注册功能：

