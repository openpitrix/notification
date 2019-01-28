/*==============================================================*/
/* Table: address                                               */
/*==============================================================*/
create table address
(
   address_id           varchar(50)  not null,
   address_list_id      varchar(50)  not null,
   address              varchar(200) not null,
   remarks              varchar(200) null,
   verification_code    varchar(50)  null,
   create_time          timestamp    not null default CURRENT_TIMESTAMP,
   verify_time          timestamp    not null default CURRENT_TIMESTAMP,
   status_time          timestamp    not null default CURRENT_TIMESTAMP,
   notify_type          varchar(50)  not null comment 'web / mobile / email / sms',
   status               varchar(50)  not null,
   PRIMARY KEY (address_id)
);
CREATE INDEX address_address_list_id_idx ON address (address_list_id);
CREATE INDEX address_address_idx ON address (address);
CREATE INDEX address_notify_type_idx ON address (notify_type);
CREATE INDEX address_status_idx ON address (status);

/*==============================================================*/
/* Table: address_list                                          */
/*==============================================================*/
create table address_list
(
   address_list_id      varchar(50)  not null,
   name                 varchar(100) null,
   extra                json         null,
   status               varchar(50)  not null,
   create_time          timestamp    not null default CURRENT_TIMESTAMP,
   status_time          timestamp    not null default CURRENT_TIMESTAMP,
   PRIMARY KEY (address_list_id)
);
CREATE INDEX address_list_status_idx ON address_list (status);

/*==============================================================*/
/* Table: notification                                          */
/*==============================================================*/
create table notification
(
   notification_id      varchar(50)  not null,
   content_type         varchar(50)  not null default '' comment 'invite, verify, fee, business',
   title                varchar(255) not null,
   content              text         not null,
   short_content        text         not null comment 'used by sms, mobile',
   expired_days         int          not null default 0 comment 'expired days,  0 is for never',
   address_info         json         not null,
   owner                varchar(50)  not null,
   status               varchar(50)  not null comment 'pending / sending / successful / failed',
   create_time          timestamp    not null default CURRENT_TIMESTAMP,
   status_time          timestamp    not null default CURRENT_TIMESTAMP,
   PRIMARY KEY (notification_id)
);
CREATE INDEX notification_content_type_idx ON notification (content_type);
CREATE INDEX notification_owner_idx ON notification (owner);
CREATE INDEX notification_status_idx ON notification (status);

/*==============================================================*/
/* Table: nf_address_list                         */
/*==============================================================*/
create table nf_address_list
(
   nf_address_list_id   varchar(50)  not null,
   notification_id      varchar(50)  not null,
   address_list_id      varchar(50)  not null,
   PRIMARY KEY (nf_address_list_id)
);
CREATE INDEX nf_address_list_notification_id_idx ON nf_address_list (notification_id);
CREATE INDEX nf_address_list_address_list_id_idx ON nf_address_list (address_list_id);

/*==============================================================*/
/* Table: task                                                  */
/*==============================================================*/
create table task
(
   task_id              varchar(50)  not null,
   notification_id      varchar(50)  not null,
   error_code           int          not null default 0,
   status               varchar(50)  not null comment 'pending / sending / successful / failed',
   create_time          timestamp    not null default CURRENT_TIMESTAMP,
   status_time          timestamp    not null default CURRENT_TIMESTAMP,
   directive            json         not null,
   PRIMARY KEY (task_id)
);
CREATE INDEX task_notification_id_idx ON task (notification_id);
CREATE INDEX task_status_idx ON task (status);
CREATE INDEX task_error_code_idx ON task (error_code);
