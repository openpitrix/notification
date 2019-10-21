alter table email_config change email smtp_user_name varchar(100) default null;

alter table email_config add from_email_addr varchar(100)  default null;

update email_config set from_email_addr=smtp_user_name;
