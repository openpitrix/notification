ALTER TABLE notification
	ADD COLUMN available_start_time varchar(8) default '';


ALTER TABLE notification
	ADD COLUMN available_end_time varchar(8) default '';
