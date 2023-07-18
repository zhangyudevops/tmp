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
                         `name` varchar(60) COLLATE utf8mb4_general_ci NOT NULL COMMENT '应用名称',
                         `tag` varchar(100) COLLATE utf8mb4_general_ci NOT NULL COMMENT '镜像tag',
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