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



CREATE TABLE `image` (
                         `id` int unsigned NOT NULL AUTO_INCREMENT,
                         `name` varchar(60)  NOT NULL COMMENT '应用名称',
                         `tag` varchar(100)  NOT NULL COMMENT '镜像tag',
                         `status` varchar(1) NOT NULL DEFAULT '1' COMMENT '1为最新，0为历史',
                         PRIMARY KEY (`id`),
                         UNIQUE KEY `tag` (`tag`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='镜像详情';


CREATE TABLE `static` (
                          `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                          `name` varchar(50) NOT NULL COMMENT '应用名',
                          `path` varchar(255) NOT NULL COMMENT '在pod中的目录',
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存放应用以及对应pod nfs存储目录';


CREATE TABLE `upload` (
                          `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                          `name` varchar(50) NOT NULL COMMENT '升级包名',
                          `md5` varchar(50) NOT NULL COMMENT 'md5值',
                          `up_time` datetime NOT NULL COMMENT '上传时间',
                          PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='上传升级包';


CREATE TABLE `config` (
                          `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
                          `name` varchar(100) NOT NULL,
                          `value` varchar(255) DEFAULT NULL,
                          `app` varchar(50) NOT NULL,
                          `comment` text,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `name` (`name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统环境变量配置表';