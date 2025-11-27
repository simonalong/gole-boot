create database td_orm;
create database td_orm1;
create database td_orm2;

create stable td_orm.td_demo1(ts timestamp,name nchar(32),age int,address nchar(128)) tags(station nchar(128));
create stable td_orm1.td_demo1(ts timestamp,name nchar(32),age int,address nchar(128)) tags(station nchar(128));
create stable td_orm2.td_demo1(ts timestamp,name nchar(32),age int,address nchar(128)) tags(station nchar(128));

create table td_orm.td_china using td_orm.td_demo1 tags('china');
create table td_orm1.td_china using td_orm1.td_demo1 tags('china');
create table td_orm2.td_china using td_orm2.td_demo1 tags('china');

insert into td_orm.td_china values(now, 'zhou', 12, '杭州');
insert into td_orm1.td_china values(now, 'zhou', 11, '杭州');
