/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2018/12/29 16:32:23                          */
/*==============================================================*/

-- drop database notification;
-- create database notification;
-- use notification;

drop table if exists job;

drop table if exists notification_user_filter;


drop table if exists notification;

drop table if exists notification_app_user_filter;

drop table if exists task;

drop table if exists user_notification;

/*==============================================================*/
/* Table: job                                                   */
/*==============================================================*/
create table job
(
   job_id               varchar(50) not null,
   notification_id      varchar(50) not null,
   job_type             varchar(50) not null,
   addrs_str            text not null,
   job_action           varchar(50),
   exe_condition        varchar(200) comment '-- job send condition',
   total_task_count     int,
   task_succ_count      int comment '-- email job send result',
   error_code           varchar(50),
   status               varchar(50),
   created_at           timestamp not null default CURRENT_TIMESTAMP,
   updated_at           timestamp not null default CURRENT_TIMESTAMP
);

alter table job
   add primary key (job_id);

/*==============================================================*/
/* Table: notification_user_filter                               */
/*==============================================================*/
create table notification_user_filter
(
   notification_id      varchar(50),
   user_filter_type     varchar(50) not null default 'specified',
   user_filter_condition text not null,
   userids_str          text
);

/*==============================================================*/
/* Table: notification                                          */
/*==============================================================*/
create table notification
(
   notification_id      varchar(50) not null,
   content_type         varchar(50) not null default '' comment ' -- network / fee / new feature etc..  产品消息 故障通知 安全通知 财务相关 其他',
   sent_type            varchar(50) not null default '' comment 'web / mobile / email /sms  邮箱 短信 移动APP 浏览器推送',
   addrs_str            text not null,
   title                varchar(255) not null,
   content              text not null,
   short_content        text not null comment '-- used by sms, mobile',
   expired_days         int not null comment '-- expired days,  0 is for never',
   owner                varchar(50) not null,
   status               varchar(50) not null default '0' comment '-- new / sending / finished',
   created_at           timestamp not null default CURRENT_TIMESTAMP,
   updated_at           timestamp not null default CURRENT_TIMESTAMP,
   deleted_at           timestamp default NULL
);

alter table notification
   add primary key (notification_id);

/*==============================================================*/
/* Table: notification_app_user_filter                           */
/*==============================================================*/
create table notification_app_user_filter
(
   notification_id      varchar(50),
   app_id               text,
   app_versions_str     text,
   cluster_status       text
);

/*==============================================================*/
/* Table: task                                                  */
/*==============================================================*/
create table task
(
   task_id              varchar(50) not null,
   job_id               varchar(50) not null,
   email_addr           varchar(50) not null,
   task_action          varchar(50) default '',
   error_code           int,
   status               varchar(50),
   created_at           timestamp not null default CURRENT_TIMESTAMP,
   updated_at           timestamp not null default CURRENT_TIMESTAMP
);

alter table task
   add primary key (task_id);

/*==============================================================*/
/* Table: user_notification                                     */
/*==============================================================*/
create table user_notification
(
   user_notification_id varchar(50) not null,
   notification_id      varchar(50),
   user_id              varchar(50),
   status               varchar(50) comment '已读 未读',
   created_at           timestamp not null default CURRENT_TIMESTAMP,
   updated_at           timestamp not null default CURRENT_TIMESTAMP,
   deleted_at           timestamp default NULL
);

-- alter table user_notification
--    add primary key (user_notification_id);
--
-- alter table job add constraint FK_Reference_15 foreign key (notification_id)
--       references notification (notification_id) on delete restrict on update restrict;
--
-- alter table notification_user_filter add constraint FK_Reference_12 foreign key (notification_id)
--       references notification (notification_id) on delete restrict on update restrict;
--
-- alter table notification_app_user_filter add constraint FK_Reference_14 foreign key (notification_id)
--       references notification (notification_id) on delete restrict on update restrict;
--
-- alter table task add constraint FK_Reference_13 foreign key (job_id)
--       references job (job_id) on delete restrict on update restrict;
--
-- alter table user_notification add constraint FK_Reference_16 foreign key (notification_id)
--       references notification (notification_id) on delete restrict on update restrict;

