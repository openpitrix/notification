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
   updated_at           timestamp not null default CURRENT_TIMESTAMP,
   PRIMARY KEY (job_id)
);

CREATE INDEX job_notification_id_idx ON job (notification_id);
CREATE INDEX job_job_type_idx ON job (job_type);
CREATE INDEX job_job_action_idx ON job (job_action);
CREATE INDEX job_status_idx ON job (status);

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
   PRIMARY KEY (notification_id)
);

CREATE INDEX notification_sent_type_idx ON notification (sent_type);
CREATE INDEX notification_owner_idx ON notification (owner);
CREATE INDEX notification_status_idx ON notification (status);

create table task
(
   task_id              varchar(50) not null,
   job_id               varchar(50) not null,
   email_addr           varchar(50) not null,
   task_action          varchar(50) default '',
   error_code           int,
   status               varchar(50) comment ' new / sending / finished',
   created_at           timestamp not null default CURRENT_TIMESTAMP,
   updated_at           timestamp not null default CURRENT_TIMESTAMP,
   PRIMARY KEY (task_id)
);

CREATE INDEX task_job_id_idx ON task (job_id);
CREATE INDEX task_task_action_idx ON task (task_action);
CREATE INDEX task_status_idx ON task (status);
