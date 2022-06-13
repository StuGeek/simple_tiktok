# simple_tiktok

## 文件夹结构

```
simple_tiktok
|-- controller          // 控制器层，主要负责接受客户端的请求，并调用服务层的服务进行处理，返回响应
  |-- comment_test.go
  |-- comment.go        // 评论或取消评论接口、获取评论列表接口
  |-- favorite_test.go
  |-- favorite.go       // 点赞或取消点赞接口，获取喜欢列表接口
  |-- feed_test.go
  |-- feed.go           // 拉取视频接口
  |-- publish_test.go
  |-- publish.go        // 发布视频接口、获取作品列表接口
  |-- relation_test.go
  |-- relation.go       // 关注或取消关注接口，获取关注列表或粉丝列表接口
  |-- user_test.go
  |-- user.go           // 用户登录注册接口、获取用户信息接口
|-- global              // 全局数据文件夹
  |-- common.go         // 全局会使用到的一些结构体
  |-- config.go         // 全局配置
|-- public              // 用户投稿视频后保存在本地的文件夹
  |-- bear.mp4          // 本地视频文件，熊的视频
|-- repository          // 存储层，主要负责将数据保存在合适的结构体里并与数据库进行交互
  |-- comment_dao.go    // 包含评论相关结构体和数据库方法
  |-- db_init.go        // 初始化数据库
  |-- demo_data.go      // demo数据，包括导入方法
  |-- favorite_dao.go   // 包含点赞相关结构体和数据库方法
  |-- relation_dao.go   // 包含关注相关结构体和数据库方法
  |-- user_dao.go       // 包含用户信息相关结构体和数据库方法
  |-- video_dao.go      // 包含视频相关结构体和数据库方法
|-- service             // 服务层，主要负责调用存储层相关方法为控制器层进行相关逻辑处理，返回控制器层所需结果
  |-- comment.go        // 评论功能相关
  |-- favorite.go       // 点赞功能相关
  |-- feed.go           // 拉取视频功能相关
  |-- publish.go        // 投稿视频功能相关
  |-- relation.go       // 关注功能相关
  |-- user.go           // 登录注册、获取用户信息相关
|-- main.go
|-- router.go
```

## 运行方式

### 运行环境

+ 服务端：Windows\Linux\Mac
+ 客户端：安卓虚拟机\真机
+ Go版本：Go 1.17+
+ 数据库：mysql 8.0+

### 运行步骤

#### 1. 设置服务器url和数据库相关信息

首先在`global/config.go`全局配置文件中，修改变量`ServerUrl`的值，设置服务器的url，如`var ServerUrl = "http://127.0.0.1:8080/"`，设置连接数据库的用户名`SqlUsername`、密码`SqlPassword`、数据库名`SqlDBName`，如果需要直接导入demo数据的话，设置导入demo数据使用的数据库名`SqlDemoDBName`（使用Demo数据会先清空之前的数据库再导入数据，如果想使用之前程序保存在数据库中的数据，需要使用名`SqlDBName`的普通数据库）

#### 2. 启动服务端

在终端进入`simple_tiktok`文件目录下，输入

```
go run main.go router.go
```

如果需要使用demo数据库并导入demo数据，输入

```
go run main.go router.go --demo
```

即可在默认的8080端口运行起来服务端

#### 3. 启动客户端

在安卓虚拟机或真机中安装`app-release.apk`，然后在客户端打开软件，连续点击右下方的“我”2~3次，弹出“高级设置”界面，设置服务器url，和之前在`global/config.go`文件中设置的服务器url一样，如`http://127.0.0.1:8080/`，点击保存并重启，退出软件，再点击进入，即可完成相应设置

![](./imgs/1.png)

## 功能说明

+ 登录、注册功能：每次启动程序会将用户信息导入到内存中，方便其它功能查询，注册会同时更新内存和数据库中的用户信息
+ 视频 Feed 流：支持所有用户刷抖音，包括未登录的用户，视频按照投稿时间倒序推出，单词最多30个
+ 视频投稿：要求用户必须处在登录状态，可以选择自己拍的视频上传，视频上传后会保存到本地 public 目录中，访问时用 127.0.0.1:8080/static/video_name 即可
+ 个人信息：登录用户可以在个人信息页查看自己的用户信息，包括用户名、关注数和粉丝数等
+ 点赞功能：登录用户可以对视频点赞或取消点赞，在个人主页能够查看点赞视频列表，未登录状态下也可以看到视频的总点赞数
+ 评论功能：登录用户可以对视频评论，未登录状态下可以看到视频的总评论数，也可以查看所有用户的评论，按发布时间倒序
+ 关注功能：登录用户可以关注或取消关注其他用户，可以在个人信息页点击打开关注列表和粉丝列表

## API接口文档

[API接口文档](https://www.apifox.cn/apidoc/project-1014925/api-20967126)

## 实现思路

### 整体设计

![](./imgs/8.png)

整体设计上，将程序分为控制层、服务层和存储层，控制层负责接口设计，接收来自客户端的请求，使用接收到的请求数据参数调用服务器层的方法进行处理，并最后响应结果；服务层负责具体处理业务逻辑，使用存储层的数据结构和方法与数据库进行交互，为控制层提供服务；存储层负责使用合适的数据结构存储信息，并提供与数据库进行交互的方法，使服务层可以利用；使用的存储数据库是Mysql。

### 数据库设计

数据库使用的数据结构和交互相关的代码文件放在`repository`文件夹中，设计了5个表，分别是用户信息表`users`、视频信息表`videos`、评论信息表`comments`、点赞视频信息表`favorites`、关系信息表`relations`。

#### 1. 用户信息表users

用户信息表users使用到的用户信息结构体`UserDao`仿照程序使用的`global/common.go`文件中的用户结构体`User`进行设计，放在`repository/user_dao.go`文件中，设置用户ID为自增主键，密码经过加密存储在数据库中，token按照一定的规则生成，有一个过期时间，过了这个时间进行某些需要登录状态下的操作时会进行检验，重新进行登录才能继续进行操作：

```go
// 用户信息表users
type UserDao struct {
	Id                int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name              string `json:"name"`
	FollowCount       int64  `json:"follow_count" gorm:"default:0"`
	FollowerCount     int64  `json:"follower_count" gorm:"default:0"`
	IsFollow          bool   `json:"is_follow" gorm:"default:false"`
	Password          string `json:"password"`
	Token             string `json:"token"`
	TokenLastUsedTime int64  `json:"token_last_used_time"`
}


func (UserDao) TableName() string {
	return "users"
}
```

#### 2. 视频信息videos表

视频信息表videos使用到的视频结构体`VideoDao`仿照程序使用的`global/common.go`文件中的视频结构体`Video`进行设计，放在`repository/video_dao.go`文件中，`Video`的`User`类型字段`Author`改成了`int64`类型的字段`AuthorId`，在数据库中只存储视频作者的用户Id，当要获取这个对应作者用户的信息时，只需要通过这个表项与`users`表进行内连接，即可获取详细的用户信息，同时新加了一个`PublishTime`字段用来表示投稿时间，放在`repository/video_dao.go`文件中，设置视频ID为自增主键：

```go
// 视频信息表videos
type VideoDao struct {
	Id            int64  `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT"`
	AuthorId      int64  `json:"author_id,omitempty" gorm:"index"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
	PublishTime   int64  `json:"publish_time,omitempty" gorm:"index:,sort:desc"`
}

func (VideoDao) TableName() string {
	return "videos"
}
```

#### 3. 点赞视频信息favorites表

点赞视频信息表favorites使用到的视频结构体`FavoriteDao`有两个字段，分别是给这个视频点赞的用户的Id和这个视频的Id，放在`repository/favorite_dao.go`文件中：

```go
// 点赞视频信息表favorites
type FavoriteDao struct {
	UserId  int64 `json:"user_id" gorm:"index"`  // 用户的Id
	VideoId int64 `json:"video_id" gorm:"index"` // 用户点赞的视频Id
}

func (FavoriteDao) TableName() string {
	return "favorites"
}
```

#### 4. 评论信息comments表

评论信息表comments使用到的评论结构体`CommentDao`仿照程序使用的`global/common.go`文件中的视频结构体`Comment`进行设计，`Comment`的`User`类型字段`User`改成了`int64`类型的字段`UserId`，在数据库中只存储评论者的Id，当要获取这个评论者的具体用户信息时，只需要通过这个表项与`users`表进行内连接，即可获取详细的用户信息，同时新加了`VideoId`字段表示视频，通过与`videos`表进行内连接即可获得详细视频信息，`PublishTime`字段用来表示评论时间，放在`repository/comment_dao.go`文件中，设置评论ID为自增主键：

```go
// 评论信息表comments
type CommentDao struct {
	Id          int64  `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserId      int64  `json:"user_id"`
	VideoId     int64  `json:"video_id" gorm:"index"`
	Content     string `json:"content"`
	CreateDate  string `json:"create_date"`
	PublishTime int64  `json:"publish_time" gorm:"index:,sort:desc"`
}

func (CommentDao) TableName() string {
	return "comments"
}
```

#### 5. 关系信息relations表

关注或取消关注过程中使用到的与数据库进行交互的用户信息结构体是`RelationDao`有两个字段，分别是给这个关注者的用户id和被关注者的用户id，放在`repository/relation_dao.go`文件中：

```go
// 关系信息表relations
type RelationDao struct {
	UserId   int64 `json:"user_id" gorm:"index"`    // 关注者的用户Id
	ToUserId int64 `json:"to_user_id" gorm:"index"` // 被关注者的用户Id
}

func (RelationDao) TableName() string {
	return "relations"
}
```

### 服务层功能实现思路

#### 1. 登录、注册功能

注册时，首先判断用户名或密码是否为空和超过32个字符，为空或超过则直接返回注册失败，接着判断用户是否存在，如果用户已经存在，直接返回注册失败，否则对密码进行加密，根据一定规则生成`token`，在数据库`users`表中创建相应的用户记录

登录时，首先判断用户名或密码是否为空，为空则直接返回登录失败，然后查询用户名，不存在则返回登录失败，再进行密码校验，将经过加密后的密码与数据库中存储着的密码进行对比，校验成功后，更新token的上次使用时间为当前时间，登录成功。

#### 2. 视频 Feed 流

当拉取视频时，首先获取限制返回视频的投稿时间戳`latest_time`，如果没有设置，则默认为当前时间，然后找到投稿时间不晚于lastestTime的投稿视频，按投稿时间倒序排列，最多`MaxFeedVideosNumOnce`个，如果没有直接返回nil，同时获取用户点赞的视频列表，并存储在map中，根据投稿视频是否在这个点赞的视频列表中，设置获取的投稿视频的是否点赞`IsFavorite`属性，并记录本次返回的视频中，发布最早的时间`nextTime`，作为下次请求时的`latest_time`，最后返回视频列表

#### 3. 视频投稿

当投稿时，首先获取token，判断是否处于登录状态以及登录凭证是否有效，不是则直接返回，取消发布视频，否则获取发布的视频文件数据，将视频文件经过文件路径和命名处理后存入本地的`public`文件夹目录下，可以添加用户序号和时间戳避免同名文件覆盖，然后在数据库中创建相应的视频记录。

当获取发布作品列表时，首先从数据库中的`videos`表中获取登录用户的点赞列表和查看的用户的发布视频列表，根据登录用户是否点赞设置视频列表中每个视频的`IsFavorite`属性，最后返回查看的那一个用户发布的视频列表。

#### 4. 个人信息

在登录时会获取个人信息，首先会根据token获取用户，如果token对应的用户不存在，返回用户不存在，如果用户存在但是token已经过期，返回token过期，需要重新登录，如果token没过期那就更新token的上次使用时间，最后返回具体用户信息。

#### 5. 点赞功能

当点赞或取消点赞时，首先获取token，判断是否处于登录状态以及登录凭证是否有效，不是则直接返回，否则先根据token获取用户的Id，然后根据Id获取这个用户点赞的视频列表，判断用户是否对当前视频点过赞了，如果是点赞行为且之前没有给这个视频点过赞，那么将数据库中`videos`表的这个视频的总点赞数加一，并在`favorites`点赞视频表中创建相应点赞记录，之前点过赞则不作反应，直接返回；如果是取消点赞行为且之前给这个视频点过赞了，那个更新数据库中的视频总点赞数，删除点赞记录，之前没点过赞则不作反应，直接返回。

当获取点赞列表时，可以根据用户Id，用`favorite_videos`表的`video_id`字段和`videos`表的`id`字段内连接，以及`videos`表的`author_id`字段和`users`表的`id`字段进行内连接查询出特定Id对应用户所点赞的视频和视频的作者，最后返回。

#### 6. 评论功能

当进行评论或取消评论时，首先获取token，判断是否处于登录状态以及登录凭证是否有效，不是则直接返回，如果是评论行为，首先获取当前日期作为创建日期，获取当前时间作为发布时间，然后向评论信息表中插入相应的评论记录，并更新视频信息表中相应视频的评论数加一；如果是取消评论行为，则从评论信息表中删除相应的记录，并更新视频信息表中相应视频的评论数减一。

当获取视频的所有评论时，可以根据视频的Id，用`comments`表的`user_id`字段和`users`表的`id`字段进行内连接查询出特定视频Id的评论和对应的评论用户，返回相应的评论列表。

#### 7. 关注功能

当进行关注或取消关注时，首先获取token，判断是否处于登录状态以及登录凭证是否有效，不是则直接返回，否则先根据token获取这个用户的Id，如果是用户自己关注或取关自己，则不能操作，直接返回，否则继续，如果是关注行为，那么在数据库的`relations`表中创建相应的记录，在`users`表中更新关注用户和被关注用户的关注数和被关注数；如果是取消关注行为，那么在数据库的`relations`表中删除相应的记录，在`users`表中更新关注用户和被关注用户的关注数和被关注数。

当获取关注列表或粉丝列表时，首先从数据库中的`relations`表中获取登录用户的关注列表和查看的用户的关注列表或粉丝列表，根据登录用户是否关注设置查看用户返回的关注列表或粉丝列表中，每个用户的`IsFollow`属性，最后返回查看的那一个用户发布的关注列表或粉丝列表。

## 功能展示

[登录、注册功能演示视频](http://120.79.66.18:8080/static/videos/%E7%99%BB%E5%BD%95%E3%80%81%E6%B3%A8%E5%86%8C%E5%8A%9F%E8%83%BD.mp4)

[拉取视频流功能演示视频](http://120.79.66.18:8080/static/videos/%E6%8B%89%E5%8F%96%E8%A7%86%E9%A2%91%E5%8A%9F%E8%83%BD.mp4)

[投稿功能演示视频](http://120.79.66.18:8080/static/videos/%E6%8A%95%E7%A8%BF%E5%8A%9F%E8%83%BD.mp4)

[获取作品列表和喜欢列表演示视频](http://120.79.66.18:8080/static/videos/%E8%8E%B7%E5%8F%96%E4%BD%9C%E5%93%81%E5%88%97%E8%A1%A8.mp4)

[点赞功能演示视频](http://120.79.66.18:8080/static/videos/%E7%82%B9%E8%B5%9E%E5%8A%9F%E8%83%BD.mp4)

[评论功能演示视频](http://120.79.66.18:8080/static/videos/%E8%AF%84%E8%AE%BA%E5%8A%9F%E8%83%BD.mp4)

[关注功能演示视频](http://120.79.66.18:8080/static/videos/%E5%85%B3%E6%B3%A8%E5%8A%9F%E8%83%BD.mp4)

这里导入Demo数据进行展示，Demo数据写在`controller/demo_data.go`文件中，使用数据库`demo_simple_tiktok`存放导入的Demo数据，在服务器终端输入`go run main.go router.go --demo`即可导入Demo数据，在Demo数据中：

+ 有5个用户user1、user2、user3、user4、user5，用户名和密码分别为：
  + 用户名：user1，密码：111111
  + 用户名：user2，密码：222222
  + 用户名：user3，密码：333333
  + 用户名：user4，密码：444444
  + 用户名：user5，密码：555555
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

### 1. 登录、注册功能：

![请添加图片描述](https://img-blog.csdnimg.cn/fe3db144ce134ef2b8957f176949a086.gif)

一开始在Demo数据库中，没有test1这个用户，所以当进行登录时，显示用户不存在，登录失败，然后进行注册，使用用户名`test1`和密码`123456`注册成功之后，退出登录，回到登录页面后，以用户名`test1`和密码`1234567`进行登录时，因为密码错误，所以登录失败，将密码改变成`123456`之后，登录成功

### 2. 视频 Feed 流

![请添加图片描述](https://img-blog.csdnimg.cn/54b5e4101fba4535acd7250a4b1499ff.gif)

在无登录状态下，打开软件，拉取视频流，因为Demo数据库中的视频一共有31个，Id由1到31，投稿时间与Id相同，在每个视频左下角的内容中有每个视频的序号和发布时间，所以应该拉取到的视频是发布时间从31到2的第31号到第2号视频，可以看到，进入软件后，第一个视频就是第31号视频（甜甜圈视频，左下角显示是第31号视频，发布时间为31），之后继续向下划动屏幕，可以看到，依次是第30号视频和第29号视频，往上划动又是第30号和第31号视频，一直往下划动，直到最后一个发布时间为2的第2号视频（橙子视频），可以看到，继续向下划动，出现的仍是第2号视频，说明这批视频中第2号视频是最后一个视频，刚好30个，符合一次最多返回30个视频，按投稿时间倒序排列的要求，在最后一个视频继续向下划动，会重新拉取视频流，以上次返回的视频中，发布最早的时间，作为下次请求时的latest_time，此时返回的应该是第2号视频和第1号视频（熊视频），所以之前第2号视频继续向下划动，依然是第2号视频，此时作为下一批视频的开头。

### 3. 视频投稿

![请添加图片描述](https://img-blog.csdnimg.cn/e9c478022d08470a87dda2d520b58ad9.gif)

用户test1在登录状态下，进行投稿，选择了本地的椰子视频，视频内容为good，投稿完成后，重新打开软件拉取视频流，可以看到，第一个视频就是刚刚投稿的视频，将服务器停止后，重新打开服务器，要在`global/config.go`中设置使用的数据库名`SqlDBName`为Demo数据库`demo_simple_tiktok`，然后在终端输入`go run main.go router.go`，重新启动服务器，再打开软件，看到第一个视频仍是刚刚投稿的视频，说明视频保存在了本地，路径保存在了数据库中，实现了持久化存储。

### 4. 个人信息

![请添加图片描述](https://img-blog.csdnimg.cn/a803a915ace84029b9a1fc18bb570ed3.gif)

点击右下方的我，可以看到登录用户的关注、粉丝、作品、喜欢等情况，在视频页面向右划动，则可以看到每个视频作者的个人信息，包括关注、粉丝、作品、喜欢等。

### 5. 点赞功能

![请添加图片描述](https://img-blog.csdnimg.cn/0c96137aebc94b3aa976a795c5abbebd.gif)

首先登录user1用户，根据Demo数据，user1没有给Id为31的甜甜圈视频点赞，但给了下一个Id为30的视频点赞，可以看到，user1在甜甜圈视频的点赞按钮并没有被点亮，点赞后，user5的个人界面中显示有一获赞，而user1个人页面中的喜欢列表中也出现了甜甜圈视频，再给甜甜圈视频取消点赞，可以看到user5的个人界面中显示去掉获赞，而user1个人页面中的喜欢列表中甜甜圈视频也消失了，再给甜甜圈视频点赞，此时甜甜圈视频的点赞数为4，往下划到第30个视频，user1之前给这个视频点过赞，所以点赞按钮是亮着的，退出登录user1用户后重新登录user5用户，可以看到此时甜甜圈视频的点赞数已经变成了4，由于客户端问题，本来user5应该是给甜甜圈视频点过赞的，但是点赞按钮没亮，此时点击点赞按钮，服务器会阻止点赞行为，显示之前已经点过赞了，向下划动到第30号视频，给其点赞，可以看到user5的喜欢列表中也出现了这视频。

### 6. 评论功能

![请添加图片描述](https://img-blog.csdnimg.cn/1a78b7f92b6149d98d4f9af227e805b8.gif)

首先登录user5用户，点开甜甜圈视频的评论列表，可以看到，获取的评论按照时间倒序排列，接着给甜甜圈视频进行评论，评论成功后，由于评论时间是最新的，所以评论会去到评论列表的最上方，长按取消评论，可以看到，取消评论功能也是正常，再评论一次，内容为`hey`，然后退出并重新登录user1用户，可以看到，甜甜圈视频的评论数变成了6，点开甜甜圈视频的评论列表，可以看到user5的评论也在评论列表里，user1也可以对甜甜圈视频进行评论。

### 7. 关注功能

![请添加图片描述](https://img-blog.csdnimg.cn/0e8db04a65d5481b94f9a8787a7ba3c8.gif)

首先登录user1，根据Demo数据，user1关注了user2，而user4、user5关注了user1，可以看到个人界面，user1的关注数为1，粉丝数为2，符合Demo数据，点开关注列表后，看到关注和粉丝的用户以及关注状态也和Demo数据相同，然后点击关注user5，可以看到，关注列表多了user5，粉丝列表中user5的关注状态也为已经关注，再关注user4，可以看到，关注列表和关注状态也做了相应的更新，符合要求，取关user4，功能也是正常，然后退出并重新登录user5用户，因为刚刚被user1关注，可以看到，user5的粉丝列表中也出现了user1，关注状态也为已经关注。

## 安全问题考虑

### 数据库注入问题

为了防止数据库注入，这里使用了gorm框架的参数化查询方法进行避免，而不是直接拼接sql语句，gorm框架的参数化查询会对sql语句进行预编译，而不是直接进行拼接，将用户输入的值用`?`占位符代替，之后运行时再传入用户输入的数据，除了防止sql注入以外，还可以对预编译的sql语句进行缓存，运行时就省去了解析优化sql语句的过程，可以加速sql的查询，提高查询性能。
