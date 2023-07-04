CREATE TABLE `user` (
                        `id` int(8) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增id',
                        `username` varchar(30) NOT NULL COMMENT '用户名',
                        `password` varchar(100) NOT NULL COMMENT '密码',
                        `nickname` varchar(40) DEFAULT NULL COMMENT '用户别名',
                        `email` varchar(40) DEFAULT NULL COMMENT '邮箱',
                        `salt` varchar(40) NOT NULL COMMENT '密码加盐',
                        PRIMARY KEY (`id`),
                        UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';


CREATE TABLE `setting` (
                           `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                           `config` text NOT NULL COMMENT '配置文件变量全量json',
                           PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配置变量';