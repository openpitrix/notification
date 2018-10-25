https://blog.csdn.net/silence_stone/article/details/51243430

#1.设置FQDN
输入下面的命令查看当前的主机名。 
hostname -f
如果你的Ubuntu服务器还没有设置好主机名，可以使用hostnamectl来设置。
 
hostnamectl set-hostname mail.app-center.cn 
通常邮箱服务器的FQDN主机名为mail.app-center.cn。


#2.安装Postfix    
sudo apt-get update   

apt-get install postfix  
apt-get install sasl2-bin 


#3.修改配置文件
vim /etc/postfix/main.cf

#文件末尾加上
smtpd_sasl_auth_enable = yes
smtpd_sasl_local_domain = $myhostname
smtpd_recipient_restrictions = permit_mynetworks,permit_sasl_authenticated,reject_unauth_destination
smtpd_client_restrictions = permit_sasl_authenticated
broken_sasl_auth_clients = yes
smtpd_sasl_path = smtpd
smtpd_sasl_security_options = noanonymous
smtpd_sasl_authenticated_header = yes
message_size_limit = 15728640

#/etc/postfix/sasl/目录下创建文件smtpd.conf,内容为： 
pwcheck_method: auxprop
auxprop_plugin: sasldb
mech_list: PLAIN LOGIN CRAM-MD5 DIGEST-MD5 NTLM 

使用saslpasswd2创建用户：
#saslpasswd2 -c -u test.com test
saslpasswd2 -c -u app-center.cn openpitrix

cp -a /etc/sasldb2 /var/spool/postfix/etc/ 
//这里很关键，在ubuntu下postfix所能浏览的目录有限制，
必须把数据库文件复制到postfix的运行目录下，不然在用户验证的时候会出错。

将postfix添加到sasl组： 
gpasswd -a postfix sasl

修改sasldb权限 
chmod 640 /var/spool/postfix/etc/sasldb2

列举sasldb2中的用户 
sasldblistusers2 -f /var/spool/postfix/etc/sasldb2

#telnet验证： 

telnet mail.app-center.cn 25
ehlo hello ##链接成功  
auth login

mail from: openpitrix@app-center.cn ##设置邮件发送端 
250 2.1.0 Ok  
rcpt to: huojiao2006@163.com  ##设置接受端 
data
ssss
.  

然后输入邮箱账户的base64码，在输入密码的base64码。
http://tool.oschina.net/encrypt?type=3

openpitrix@app-center.cn  
openpitrix
 
b3BlbnBpdHJpeEBhcHAtY2VudGVyLmNu
b3BlbnBpdHJpeA==
 
 
service saslauthd status
sevice postfix status


























