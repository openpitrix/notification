
alter table address_list CHANGE name address_list_name varchar(50);
alter table address_list alter column status SET DEFAULT 'active';


alter table address drop COLUMN address_list_id;

create table address_list_binding
(
   binding_id           varchar(50) not null,
   address_list_id      varchar(50),
   address_id           varchar(50),
   create_time          timestamp default CURRENT_TIMESTAMP
);

alter table address_list_binding
   add primary key (binding_id);



alter table notification alter column short_content text   null comment 'used by sms, mobile';
