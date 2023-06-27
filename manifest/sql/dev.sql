create table user
(
    id       int auto_increment Primary Key comment '自增id',
    username varchar(30) not null unique comment '用户名',
    password varchar(100) not null comment '密码',
    nickname varchar(40) null comment '用户别名',
    email    varchar(40) null comment '邮箱',
    salt     varchar(40) not null comment '密码加盐'
)
    comment '系统用户表';
