/*==============================================================*/
/* Table: email_config                                          */
/*==============================================================*/
create table email_config
(
   protocol             varchar(30),
   email_host           varchar(100),
   port                 int,
   display_sender       varchar(100),
   email                varchar(100),
   password             varchar(100),
   ssl_enable           bool,
   create_time          timestamp not null default CURRENT_TIMESTAMP,
   status_time          timestamp not null default CURRENT_TIMESTAMP
);


INSERT INTO email_config (protocol,email_host,port,display_sender,email,password,ssl_enable) VALUES
('SMTP','mail.app-center.com.cn',25,'notification','admin@app-center.com.cn','password',0);

