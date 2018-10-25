CREATE TABLE email_job
(
    email_job_id varchar(50) PRIMARY KEY NOT NULL,
    title text DEFAULT '' NOT NULL,                              -- email title
    condition text DEFAULT '' NOT NULL,                          -- email job send condition
    content text DEFAULT '' NOT NULL,                            -- email content
    succ_count integer DEFAULT 0 NOT NULL,                       -- email succedded send count
    total_count integer DEFAULT 0 NOT NULL,                      -- email total send count
    result text DEFAULT '' NOT NULL,                             -- email job send result
    create_time timestamp DEFAULT current_timestamp NOT NULL,
    console_id varchar(50) DEFAULT '' NOT NULL,
    status varchar(16) NOT NULL,
    status_time timestamp DEFAULT current_timestamp NOT NULL,
    owner varchar(255) NOT NULL
);
CREATE INDEX email_job_create_time_index ON email_job (create_time);
CREATE INDEX email_job_owner_index ON email_job (owner);
CREATE INDEX email_job_status_index ON email_job (status);

-- sms_job table
CREATE TABLE sms_job
(
    sms_job_id varchar(50) PRIMARY KEY NOT NULL,
    sms_type varchar(50) NOT NULL,                  -- register, change_phone, balance_almost_depleted, insufficient_balance_notify
    phone varchar(20) NOT NULL,
    content varchar(255) NOT NULL,
    status varchar(16) NOT NULL,
    create_time timestamp DEFAULT current_timestamp NOT NULL,
    tick_id varchar(50) NOT NULL,               -- tick id from sms provider
    provider varchar(16) NOT NULL
);
CREATE INDEX sms_job_status_index ON sms_job (status);
CREATE INDEX sms_job_type_index ON sms_job (sms_type);
CREATE INDEX sms_job_phone_index ON sms_job (phone);


--notification_center_post 是通知内容，包括通知的方式，通知类型 

-- notification_center_post
CREATE TABLE notification_center_post
(
    notification_center_post_id varchar(50) PRIMARY KEY NOT NULL,
    post_type varchar(255) DEFAULT '' NOT NULL,                      -- network / fee / new feature etc..
    user_filter_type varchar(50) DEFAULT 'specified' NOT NULL,      -- all / specified / condition
    user_condition text,                                            -- user condition
    specified_users text,                                           -- user list
    notify_type varchar(255) DEFAULT '' NOT NULL,                    -- web / mobile / email /sms
    title varchar(1024) NOT NULL,
    content text NOT NULL,
    short_content text NOT NULL,                                    -- used by sms, mobile
    owner varchar(255) NOT NULL,
    console_id varchar(50) DEFAULT '' NOT NULL,
    status varchar(50) NOT NULL,                                    -- new / sending / finished
    expired_days integer DEFAULT 0 NOT NULL,                        -- expired days,  0 is for never
    create_time timestamp DEFAULT current_timestamp NOT NULL,
    status_time timestamp DEFAULT current_timestamp NOT NULL
);
CREATE INDEX notification_center_post_post_type_index ON notification_center_post(post_type);
CREATE INDEX notification_center_post_user_filter_type_index ON notification_center_post(user_filter_type);
CREATE INDEX notification_center_post_status_index ON notification_center_post(status);
CREATE INDEX notification_center_post_create_time_index ON notification_center_post(create_time);
CREATE INDEX notification_center_post_status_time_index ON notification_center_post(status_time);

CREATE TABLE notification_center_post_user
(
    post_id varchar(20) NOT NULL,
    user_id varchar(20) NOT NULL,
    status varchar(10) NOT NULL,                                    -- new / read
    create_time timestamp DEFAULT current_timestamp NOT NULL
);
CREATE INDEX notification_center_post_user_post_index ON notification_center_post_user(post_id);
CREATE INDEX notification_center_post_user_user_index ON notification_center_post_user(user_id);
CREATE INDEX notification_center_post_user_status_index ON notification_center_post_user(status);
CREATE INDEX notification_center_post_user_create_time_index ON notification_center_post_user(create_time);

--notification_map 设置某种通知类型 对应的通知列表，比如 财务 设置对应的通知列表
-- user nf map
CREATE TABLE notification_map
(
    notification_type varchar(50) NOT NULL,
    notification_list_id varchar(50) NOT NULL,
    owner varchar(255) NOT NULL,
    root_user_id varchar(255) NOT NULL,
    console_id varchar(50) NOT NULL,
    create_time timestamp DEFAULT current_timestamp NOT NULL,
    status_time timestamp DEFAULT current_timestamp NOT NULL,
    unsubscribe integer DEFAULT 0 NOT NULL,
    PRIMARY KEY(owner, notification_type)
);
CREATE INDEX notification_map_owner_index ON notification_map(owner);
CREATE INDEX notification_map_root_user_id_index ON notification_map(root_user_id);
CREATE INDEX notification_map_console_id_index ON notification_map(console_id);



-- nf list
--notification_list 就是通知列表
CREATE TABLE notification_list
(
    notification_list_id varchar(50) PRIMARY KEY NOT NULL,
    notification_list_name text DEFAULT '' NOT NULL,
    create_time timestamp DEFAULT current_timestamp NOT NULL,
    console_id varchar(50) DEFAULT '' NOT NULL,
    owner varchar(255) NOT NULL,
    root_user_id varchar(255) DEFAULT '' NOT NULL
);
CREATE INDEX notification_list_owner_index ON notification_list (owner);
CREATE INDEX notification_list_console_index ON notification_list  (console_id);
CREATE INDEX notification_list_root_user_id_index ON notification_list (root_user_id);

--notification_item 是通知列表的联系方式
--notification_item 里面的content是 123423423423 或者 xxxx@qq.com，这些联系方式
-- nf item
CREATE TABLE notification_item
(
    notification_item_id varchar(50) PRIMARY KEY NOT NULL,
    notification_item_type varchar(50) NOT NULL, -- email, phone
    content varchar(255) NOT NULL,
    remarks text DEFAULT '' NOT NULL,
    verification_code varchar(255) NOT NULL, -- '' means verified.
    create_time timestamp DEFAULT current_timestamp NOT NULL,
    verify_time timestamp DEFAULT current_timestamp NOT NULL,
    console_id varchar(50) DEFAULT '' NOT NULL,
    owner varchar(255) NOT NULL,
    root_user_id varchar(255) DEFAULT '' NOT NULL
);
CREATE INDEX notification_item_owner_index ON notification_item (owner);
CREATE INDEX notification_item_console_index ON notification_item  (console_id);
CREATE INDEX notification_item_root_user_id_index ON notification_item (root_user_id);
CREATE INDEX notification_item_content_index ON notification_item (content);


--notification_list_item 是通知列表跟联系方式的映射关系
-- nf list item
CREATE TABLE notification_list_item
(
    notification_list_id varchar(50) NOT NULL,
    notification_item_id varchar(50) NOT NULL,
    PRIMARY KEY(notification_list_id, notification_item_id)
);


-- cluster nf
CREATE TABLE cluster_notification
(
    cluster_notification_id varchar(50) PRIMARY KEY NOT NULL,
    zones text,
    app_ids text,
    app_versions text,
    app_users text,
    cluster_status text,
    title varchar(1024) NOT NULL,
    content text NOT NULL,
    recipients text NOT NULL,
    notification_center_post_id varchar(50) NOT NULL,
    owner varchar(255) NOT NULL,
    console_id varchar(50) DEFAULT '' NOT NULL,
    status varchar(50) NOT NULL,
    create_time timestamp DEFAULT current_timestamp NOT NULL,
    status_time timestamp DEFAULT current_timestamp NOT NULL
);
CREATE INDEX cluster_notification_notification_center_post_id_index ON cluster_notification(notification_center_post_id);
CREATE INDEX cluster_notification_owner_index ON cluster_notification(owner);
CREATE INDEX cluster_notification_status_index ON cluster_notification(status);
CREATE INDEX cluster_notification_create_time_index ON cluster_notification(create_time);
CREATE INDEX cluster_notification_status_time_index ON cluster_notification(status_time);

