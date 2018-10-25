SET SESSION FOREIGN_KEY_CHECKS=0;

/* Drop Tables */

DROP TABLE IF EXISTS task;
DROP TABLE IF EXISTS job;
DROP TABLE IF EXISTS notfication_user_filter;
DROP TABLE IF EXISTS notification_appuser_filter;
DROP TABLE IF EXISTS notification_center_post_user;
DROP TABLE IF EXISTS notification_center_post;




/* Create Tables */

CREATE TABLE job
(
	job_id varbinary(50) NOT NULL,
	nf_post_id varchar(50) NOT NULL,
	job_type varchar(50),
	addrs_str text,
	job_action varchar(50),
	-- -- email job send condition
	exe_condition varchar(200) DEFAULT '' NOT NULL COMMENT '-- email job send condition',
	total_task_count int(11),
	-- -- email succedded send count
	task_succ_count int DEFAULT 0 NOT NULL COMMENT '-- email succedded send count',
	-- -- email job send result
	result varchar(50) DEFAULT '' NOT NULL COMMENT '-- email job send result',
	error_code int(11) DEFAULT 0 NOT NULL,
	status varchar(50) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (job_id)
);


CREATE TABLE notfication_user_filter
(
	nf_post_id varchar(50) NOT NULL,
	user_filter_type varchar(50) DEFAULT 'specified' NOT NULL,
	user_filter_condition text NOT NULL,
	userids_str text
);


CREATE TABLE notification_appuser_filter
(
	nf_post_id varchar(50) NOT NULL,
	app_id text,
	app_versions_str text,
	cluster_status text
);


CREATE TABLE notification_center_post
(
	nf_post_id varchar(50) NOT NULL,
	--  -- network / fee / new feature etc.. 
	-- 产品消息 故障通知 安全通知 财务相关 其他
	nf_post_type varchar(50) DEFAULT '' NOT NULL COMMENT ' -- network / fee / new feature etc.. 
产品消息 故障通知 安全通知 财务相关 其他',
	-- web / mobile / email /sms
	--  邮箱 短信 移动APP 浏览器推送
	notify_type varchar(50) DEFAULT '' NOT NULL COMMENT 'web / mobile / email /sms
 邮箱 短信 移动APP 浏览器推送',
	addrs_str text,
	title varchar(1024) NOT NULL,
	content text NOT NULL,
	-- -- used by sms, mobile
	short_content text NOT NULL COMMENT '-- used by sms, mobile',
	-- -- expired days,  0 is for never
	expired_days int DEFAULT 0 NOT NULL COMMENT '-- expired days,  0 is for never',
	owner varchar(50) NOT NULL,
	-- -- new / sending / finished
	status varchar(50) NOT NULL COMMENT '-- new / sending / finished',
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (nf_post_id)
);


CREATE TABLE notification_center_post_user
(
	nf_post_user_id varchar(50) NOT NULL,
	nf_post_id varchar(50) NOT NULL,
	user_id varchar(50),
	-- 已读 未读
	status varchar(50) COMMENT '已读 未读',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (nf_post_user_id)
);


CREATE TABLE task
(
	task_id varchar(50) NOT NULL,
	job_id varbinary(50) NOT NULL,
	addrs_str text,
	task_action varchar(50) DEFAULT '' NOT NULL,
	result varchar(50),
	error_code int(11),
	status varchar(50),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
	PRIMARY KEY (task_id)
);



/* Create Foreign Keys */

ALTER TABLE task
	ADD FOREIGN KEY (job_id)
	REFERENCES job (job_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE job
	ADD FOREIGN KEY (nf_post_id)
	REFERENCES notification_center_post (nf_post_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE notfication_user_filter
	ADD FOREIGN KEY (nf_post_id)
	REFERENCES notification_center_post (nf_post_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE notification_appuser_filter
	ADD FOREIGN KEY (nf_post_id)
	REFERENCES notification_center_post (nf_post_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;


ALTER TABLE notification_center_post_user
	ADD FOREIGN KEY (nf_post_id)
	REFERENCES notification_center_post (nf_post_id)
	ON UPDATE RESTRICT
	ON DELETE RESTRICT
;



