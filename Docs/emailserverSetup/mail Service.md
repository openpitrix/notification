#==========================================================================================================================================
#为什么需要搭建自己的邮件服务器？
https://www.cnblogs.com/hgj123/p/6186400.html
对于用户自己的网站来说，发送各种例如注册通知的邮件是很基本的一个需求，之前我一直用的是腾讯的企业邮箱，感觉挺方便的，直接可以绑定QQ邮箱接收邮件，网站配置一下SMTP也就可以发出邮件。但是在前几天由于有重要信息需要立即通知用户，所以选择了群发邮件的方式。在当我以为一切都是辣么完美的时候，陆续有用户过来问我什么情况，我都会跟他们说请查收邮件，但是有好几个人说并没有任何邮件，于是我试着再发一次，结果返回了错误提示。在网上找了下原因，后来看到这个：
各大免费邮箱邮件群发账户SMTP服务器配置及SMTP发送量限制情况，才知道是因为发信数量限制了。

所以我只好另寻出路了，然后我在知乎上面找到了很多个提供邮件发送的服务商，大概有这些：SendGrid、MailChimp、Amazon SES、SendCloud、Mailgun等等，在看了不少人的建议之后，我选择了Mailgun。
Mailgun注册和配置都挺简单，很快我就成功的发出了第一封邮件，怀着这封欣喜，我又发送了几封邮件，可是悲剧发生在第三封邮件，Mailgun后台有详细的发送记录，这个非常不错，在后台我看到我的邮件被拒收了，
原因大概是该服务器IP的发信频率超过腾讯邮箱限制。
所以这里就涉及到IP的问题，目前第三方的邮件发送服务普遍都是共享IP(后面还试过SendCloud、)，而共享IP并不能确定是否已经达到接收方的数量限制，一旦达到了就无法再发送。
这就是说还需要使用独立IP才能保证邮件有较高的到达率，接着就看了各家的独立IP价格，一般都是二十几甚至四十几美刀一个月，这对于我们这种小站长邮件需求不高的来说确实有点贵，买台VPS都不用这个价吧。

经过上面这些折腾，也算明白了如果要想顺畅的发出邮件的话，除了花钱，就只有自己搭建一个邮件服务器了。

如何在Ubuntu 16.04上安装并配置Postfix作为只发送SMTP服务器
https://blog.csdn.net/zstack_org/article/details/69525954
#========================================================================================================================================== 
#域名（域名（app-center.cn）是百度）
是百度云上的(上的(https://login.bce.baidu.com/)，
用户名，用户名 15827567252 密码 eam82g2,可以在这里修改DNS转向的IP


 自建邮件服务器域名解析设置(A与MX记录) 
如果域名没有做解析，只能用于内网收发邮件。
要想实现与外网邮箱的收发，需要做域名解析。是在“域名解析后台”进行设置（域名提供商提供“域名解析后台"）。

#========================================================================================
#分享几个免费的开源邮件服务器软件
分享几个免费的开源邮件服务器软件
https://blog.csdn.net/yihu0817/article/details/41966109


https://blog.csdn.net/erik_aaron/article/details/46584013

简单的替代方案——使用Postfix处理发出邮件

如果大家希望在服务器上利用简单应用发送邮件，则无需配置完整的邮件服务器。大家可以直接设置Postfix等简单邮件传输代理（简称MTA）。具体方式请参阅如何在Ubuntu 14.04中安装并设置Postfix。

另外，大家可以在自己的服务器上配置sendmail以作为发出消息的邮件传输方案。 
#========================================================================================
 #基础知识
https://blog.csdn.net/hxpjava1/article/details/80669355
1.常见邮局端口&&协议：

发邮件的协议有SMTP，收邮件的协议有POP3和IMAP。
SMTP：明文使用25端口。加密后使用25/587/465端口。
IMAP：明文使用143端口。加密后使用143/993端口。
POP3：明文使用110端口。加密后使用110/995端口。

2.常见邮局软件和安全软件：

sendmail：用于发邮件。资格最老的邮局，所有Linux发行版基本都带。但是配置麻烦。
postfix：Wietse Venema觉得sendmail配置太麻烦了，就开发了一个“简化配置版sendmail”，即postfix。支持smtp协议。
dovecot：用于收邮件，支持imap/pop3。
spamassasin：垃圾邮件过滤器。可以自订规则。
clamav：邮件杀毒工具。
opendkim：生成dkim签名。有什么用？详见下面的“反垃圾邮件技术”。
fail2ban：防止别人暴力破解用户名密码的工具。
3.反垃圾邮件技术：


运行在Linux环境下免费的邮件服务器，或者称为MTA(Mail Transfer Agent)有若干种选择，
比较常见的有Sendmail、Qmail、Postfix、exim及Zmailer等等。


首先说明基本的背景知识。一个邮件服务器通常包括如下两个基本组件：

Mail Transfer Agent (MTA)，用于向收件人的目标 agent 发送邮件和接收来自其他 agent 的邮件。我们使用 Postfix 作为 MTA，它比 sendmail 更安全高效，且在 Ubuntu 平台上官方源提供更新。
Mail Delivery Agent (MDA)，用于用户到服务器上访问自己的邮件。我们使用 Dovecot 作为 MDA，它在 Ubuntu 平台上也是官方源提供更新。

postfix 仅提供 smtp 服务，不提供 pop3 和 imap 服务，主要是用发送和接收邮件的
（接收到的邮件后，一般转交 dovecot 处理，dovecot 负责将 postfix 转发过来的邮件保存到服务器硬盘上） 

dovecot 仅提供 pop3 和 imap 服务，不提供 smtp 服务（Foxmail之类的邮箱客户端，都是通过pop3 和 imap 来收发邮件的。
发邮件时，dovecot 会将邮件转交给 postfix 来发送）

#========================================================================================
#对比方案
1.tomav/docker-mailserver
https://github.com/tomav/docker-mailserver 
利用Docker自建多功能加密邮件服务器
https://blog.csdn.net/hxpjava1/article/details/80669355


2.EwoMail
有开源版和专业版
http://www.ewomail.com/
Linux 简单搭建邮件服务器
http://gyxuehu.iteye.com/blog/2400424


3.linux安装开源邮件服务器iredmail的方法：docker
https://baike.baidu.com/item/iredmail/5314719?fr=aladdin
iRedMail 是一个基于 Linux/BSD 系统的零成本、功能完备、成熟的邮件服务器解决方案。
iRedMail 是一个开源、免费的项目。以 GPL（v2）协议发布。  
https://www.iredmail.org/  
https://www.cnblogs.com/shengulong/p/9133466.html


4.ubuntu Postfix Dovecot 
Ubuntu下Postfix邮件服务器安装及基本的设置
https://blog.csdn.net/dotuian/article/details/8552418  

使用 Postfix 设置只发送邮件的邮件服务
https://www.chrisyue.com/config-postfix-as-a-send-only-mail-server.html
 
#Ubuntu搭建简易Postfix邮箱服务器
https://blog.csdn.net/oolocal/article/details/52861583

Ubuntu搭建邮件服务器postfix
https://blog.csdn.net/MOU_IT/article/details/80256960
 
# postfix 邮件发送
https://blog.csdn.net/wolvesqun/article/details/52181172 
#======================================================================================== 
#Ubuntu搭建邮件服务器postfix  

Postfix的收发日志保存在/var/log/mail.log文件中。
Postfix本身的运行错误日志保存在/var/log/mail.err文件中
配置文件位置 /etc/postfix/main.cf
邮件位置：/var/mail/

#配置

vim /etc/postfix/main.cf

myhostname = sample.abc.com　 
设置系统的主机名 
mydomain = abc.com　 
设置域名（设置为E-mail地址“@”后面的部分） 
myorigin = $mydomain　 
发信地址“@”后面的部分设置为域名（非系统主机名）  
inet_interfaces = all　 
接受来自所有网络的请求 
mydestination = myhostname,localhost.myhostname,localhost.mydomain, localhost, $mydomain　 
指定发给本地邮件的域名 
home_mailbox = Maildir/　 
指定用户邮箱目录 
3.启动 
sudo systemctl start postfix 
4.关闭 
sudo systemctl stop postfix

postfix reload 
mailq :会列出当前在postfix发送队列中的全部邮件
postsuper -d ALL:删除当前等待发送队列的全部邮件，包含发送失败的退信

sudo hostnamectl set-hostname   mail.app-center.cn
 

echo "This is the body of the email" | mail -s "This is the subject line" huojiao2006@163.com
echo "test email from huojiao   1111" | sendmail huojiao2006@163.com

echo "test email from huojiao   1111" | sendmail johuo@yunify.com 
echo "test email" | sendmail huojiao2006@163.com 
echo "test email  from huojiao" | sendmail 13009254@qq.com  513590612@qq.com
echo "厉害了吧，哈哈哈"|mail -s "霍姣的邮件"   13009254@qq.com  513590612@qq.com

echo "厉害了吧，哈哈哈"|mail -s "霍姣的邮件"   johuo@yunify.com  huojiao2006@163.com


#chfn
这个帖子底下有一些解决方案，但是都只解决了发件人地址的问题，发件人的名字怎么都改不过来，
后来在另外一个帖子里找到了解决方案，即用chfn命令修改一下root用户下的用户信息，
postfix在发送邮件的时候就会直接去读取系统里面存储的操作账号对应的用户信息，
其他地方也有一些解决方案，但是跟用chfn大同小异，似乎并不能很方便地像我想象的那样直接在命令里面指定发件人名字。。。 
 
#postfix配置认证账号密码发送邮件
https://www.hyahm.com/article/225.html  

mail from:root@mail.app-center.cn


rcpt to:huojiao2006@163.com
RCPT TO:huojiao2006@163.com 




用过 golang 发送邮件的同学一定都知道 go 语言中默认的 smtp 模块是无法在正常的 smtp 25 端口上去发送邮件的
（有兴趣的网友可以自行用 163 的邮箱试试）。
原因是 golang 本身起源就是为 google 公司的需求服务的，所以很多功能都先优先做了 google 需要的部分，
而对电子邮件有一点了解的网友们应该都知道 google 的 gmail 是不支持常规的 smtp 25 端口的，
它需要安全连接的 ssl 接口。
所以大家如果搜索 golang 的 smtp 发送示例的话，基本上都是要进行一点改造的。
其实这样的改造代码都不完善，最后都会注明有问题。
这个问题其实源于 golang 的 smtp 源码（我看的是 1.7 版本）中对 "AUTH" 命令的实现与常规不太一样，
它的实现后面跟了两个参数，而经过我们前几篇的文章，大家都知道实现的只有一个参数，那就是 "AUTH LOGIN"。
知道了这一点，要改造 golang 的源码还是比较容易的。
不过 golang 和 java 一样有点过度设计的意思，
所以要看懂它的代码也不是太容易（不过 golang 中的各种协议代码设计得很精巧，远远不是 java 可比的）。
所以我们既然已经知道了怎样自己写一个，那还不如自己明明白白的写一个出来。










