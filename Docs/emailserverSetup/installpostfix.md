https://blog.csdn.net/oolocal/article/details/52861583
Ubuntu搭建简易Postfix邮箱服务器


#1.设置FQDN
输入下面的命令查看当前的主机名。 
hostname -f
如果你的Ubuntu服务器还没有设置好主机名，可以使用hostnamectl来设置。
 
hostnamectl set-hostname mail.app-center.cn

通常邮箱服务器的FQDN主机名为mail.yourdomain.com。

#2.安装Postfix  
更新该软件包数据库： 
sudo apt-get update  
安装mailtuils将安装Postfix以及其它一些Postfix所必需的程序。 
sudo apt install mailutils

sudo apt-get update 
sudo apt-get install postfix -y


每个用户的邮件保存在/var/spool/mail<username>和/var/mail/<username>文件中。如果你不知道收件箱保存在哪里，运行这条命令： 
postconf mail_spool_directory
Postfix的收发日志保存在/var/log/mail.log文件中。Postfix本身的运行错误日志保存在/var/log/mail.err文件中。

#3.配置Postfix  
vim /etc/postfix/main.cf
76 myhostname = westos-mail.westos.com                     ##25端口开启的网络接口 
83 mydomain = westos.com                                   ##指定mta主机名称
99 myorigin = westos.com                                   ##指定mta主机域名
116 inet_interfaces = all
164 mydestination = $myhostname, $mydomain, localhost

mynetworks=0.0.0.0/0  

**chfn命令修改root为openpitrix用户**  

**转发系统邮件** 

最后设置转发机制，这样我们就能够将指向系统root的邮件转发至自己的个人外部邮箱了。 
要实现这一功能，我们需要编辑/etc/aliases文件。 
sudo vim /etc/aliases
此文件的默认内容如下： 
/etc/aliases 
postmaster:    root 
在此设定下，系统生成的邮件会被发送至root用户。这里我们需要将其重新路由至自己的邮箱，变更后为： 
/etc/aliases 
postmaster:    root
root:          huojiao2006@163.com 
运行以下命令使变更生效： 
sudo newaliases
再次发送邮件以进行测试： 
echo "This is the body of the email"|mail -s "This is the subject line"   root


https://blog.csdn.net/Syx834722207/article/details/72667232  
默认情况下，Postfix邮件主机可以接受和转发符合以下条件的邮件：
                (1)接受邮件
                        目的地为$inet_interfaces的邮件；
                        目的地为$mydestination的邮件；
                        目的地为$vitual_alias_maps的邮件。
                (2)转发邮件
                        来自客户端IP地址符合$mynetworks的邮件；
                        来自客户端主机名称符合$relay_domains及其子域的邮件
                        目的地为$relay_domains及其子域的邮件



#4.发送测试邮件
Postfix在安装时，会同时安装一个sendmail的程序（/usr/sbin/sendmail）。
你可以用这个sendmail二进制程序向你的Gmail邮箱发送一封测试邮件。在服务器上输入下面的命令：
echo "test email" | sendmail huojiao2006@163.com

使用mail程序来发送邮件，查看收件箱 
sendmail的功能非常有限，现在让我们来安装一个命令行邮箱客户端。 
sudo apt-get install mailutils
使用mail发送邮件的命令为 
mail username@gmail.com 
echo "厉害了吧，哈哈哈"|mail -s "霍姣的邮件"   huojiao2006@163.com 

**telnet 远程发送邮件**
telnet 192.168.0.3 25
ehlo hello ##链接成功  
mail from: root@app-center.cn ##设置邮件发送端 
250 2.1.0 Ok  
rcpt to: huojiao2006@163.com  ##设置接受端 
data
ssss
.

data
354 End data with .Date: November 25, 2016
From: tester
Message-ID: first-test
Subject: mail server test
Hi carla,
Are you reading this? Let me know if you didn't get this.
.
mail from: openpitrix@app-center.cn 
rcpt to: huojiao2006@163.com 
rcpt to: johuo@yunify.com
mail from: fmaster@app-center.cn


mail from: fmaster@app-center.cn

mail from: johuo@yunify.com

#5.相关命令
service postfix start|restart|stop 

查看Postfix的版本
postconf mail_version

netstat来查看Postfix的监听情况
netstat -lnpt

Postfix的master进程监听TCP 25号端口
在发送测试邮件之前，我们最好是查看25号端口是否被防火墙或主机商屏蔽。
nmap可以帮助我们扫描服务器的开放端口。在你的个人电脑上运行下面的命令。 
sudo nmap <your-server-ip>

postfix reload
 
https://blog.csdn.net/dotuian/article/details/8552418



# 创建用户
useradd admin
#设置密码，会要求输入两次密码
passwd admin













