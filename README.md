# Agenda
Agenda是一个用于会议管理的命令行工具,每个用户都拥有一个自己的帐号
# 命令格式
```bash
agenda [Command] [arg]
```
## 用户注册
**agenda register -u username -p password -e email -t telephone**  
注册新用户时,用户需设置一个唯一的用户名(username)和一个密码(password)。另外,还需登记邮箱(email)及电话信息(telephone)。
## 用户登录
**agenda login -u username -p password**  
用户使用用户名和密码登录**Agenda**系统。
用户名为username,密码为password
## 用户登出
**agenda logout**  
已登录的用户登出系统后,只能使用用户注册和用户登录功能。
## 用户状态
**agenda status**
查看用户是否已登录
## 用户查询
**agenda showall**  
已登录的用户查看已注册的所有用户的用户名、邮箱及电话信息。
## 用户删除
**agenda deleteAccount**  
已登录的用户可以删除本用户账户（即销号）。
删除成功则退出系统登录状态。删除后,该用户账户不再存在。  
用户账户删除以后：  
以该用户为*发起者*的会议将被删除  
以该用户为*参与者*的会议将从*参与者*列表中移除该用户。若因此造成会议*参与者*人数为0,则会议也将被删除。
## 创建会议
**agenda createMeeting -t title -p participators -s startTime -e endTime**  
已登录的用户可以添加一个新会议到其议程安排中。会议可以在多个已注册用户间举行,不允许包含未注册用户。添加会议时提供的信息应包括：  
会议主题(title)（在会议列表中具有唯一性）  
会议参与者(participators),多个参与者以逗号分隔,逗号后不能有空格  
会议起始时间(start time)  
会议结束时间(end time)  
时间格式为YYYY-MM-DD-HH-mm  
如：
```bash
agenda createMeeting -t GO_lec_1 -p Tom,Mike,Shelly -s 2018-06-30-14-00 -e 2018-06-30-16-00
```
将创建一个主题为GO_lec_1的会议,会议成员有Tom, Mike, Shelly,从2018-06-30-14-00到2018-06-30-16-00的会议  
注意,任何用户都无法分身参加多个会议。如果用户已有的会议安排（作为发起者或参与者）与将要创建的会议在时间上重叠 （允许仅有端点重叠的情况）,则无法创建该会议。  
## 增删会议参与者
**agenda addParticipator -t title -p participators**  
**agenda removeParticipator -t title -p participators**  
会议主题(title)（在会议列表中具有唯一性）  
会议参与者(participators),
已登录的用户可以向自己发起的主题为**title**的会议增加/删除*参与者*（**participators**多个参与者以逗号分隔）。  
删除会议参与者后,若因此造成会议*参与者*人数为0,则会议也将被删除。  
## 查询会议
**agenda check -s startTime -e endTime**  
已登录的用户可以查询自己的议程在某一时间段(time interval)内的所有会议安排。  
用户给出所关注时间段的起始时间（**startTime**）和终止时间（**endTime**）,返回该用户议程中在指定时间范围内找到的所有会议安排的列表。如果**startTime**或**endTime**其中一项未指定,则视作没有对应的时间约束。  
在列表中给出每一会议的起始时间、终止时间、主题、以及发起者和参与者。  
注意,查询会议的结果应包括用户作为*发起者或参与者*的会议。
## 取消会议
**agenda cancel -t title**  
已登录的用户可以取消 自己发起 的某一会议安排。
取消会议时,需提供唯一标识：会议主题（title）。
## 退出会议
**agenda quit -t title**
已登录的用户可以退出自己参与的某一会议安排。  
退出会议时,需提供一个唯一标识：会议主题（title）。若因此造成会议*参与者*人数为0,则会议也将被删除。
## 清空会议
**agenda clearMeeting**  
已登录的用户可以清空*自己发起*的所有会议安排。
