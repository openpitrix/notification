 /*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2019/1/11 13:21:10                           */
/*==============================================================*/

--
-- alter table job
--    drop primary key;
--
-- drop table if exists job;
--
-- alter table notification
--    drop primary key;
--
-- drop table if exists notification;
--
-- alter table task
--    drop primary key;
--
-- drop table if exists task;

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
   status               varchar(50) comment ' new / sending / finished',
   created_at           timestamp not null default CURRENT_TIMESTAMP,
   updated_at           timestamp not null default CURRENT_TIMESTAMP
);

alter table job
   add primary key (job_id);

/*==============================================================*/
/* Table: notification                                          */
/*==============================================================*/
create table notification
(
   notification_id      varchar(50) not null,
   content_type         varchar(50) not null default '' comment ' -- network / fee / new feature etc..',
   sent_type            varchar(50) not null default '' comment 'web / mobile / email /sms',
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
/* Table: task                                                  */
/*==============================================================*/
create table task
(
   task_id              varchar(50) not null,
   job_id               varchar(50) not null,
   email_addr           varchar(50) not null,
   task_action          varchar(50) default '',
   error_code           int,
   status               varchar(50) comment ' new / sending / finished',
   created_at           timestamp not null default CURRENT_TIMESTAMP,
   updated_at           timestamp not null default CURRENT_TIMESTAMP
);

-- alter table task
--    add primary key (task_id);
--
-- alter table job add constraint FK_Reference_15 foreign key (notification_id)
--       references notification (notification_id) on delete restrict on update restrict;
--
-- alter table task add constraint FK_Reference_13 foreign key (job_id)
--       references job (job_id) on delete restrict on update restrict;
--
