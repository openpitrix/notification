ALTER TABLE task
	ADD COLUMN notify_type varchar(50) default 'email';

ALTER TABLE notification
	ADD COLUMN extra json;
