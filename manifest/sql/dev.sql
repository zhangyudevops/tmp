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



# config表初始sql
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_abcmonitor_mysql_ip', '172.27.15.154', 'abcmonitor', 'abcmonitor应用服务mysql的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_alarmengine_mysql_ip', '172.27.15.154', 'alarmengine', 'alarmengine应用服务mysql的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_abcmonitor_shentongyun_ip', '10.1.1.1', 'abcmonitor', 'abcmonitor应用服务调取深瞳云的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_abcmonitor_shentongyun_hostname', 'api.server.com', 'abcmonitor', 'sass' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_abcmonitor_replicas_number', '3', 'abcmonitor', 'abcmonitor deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_abcmonitor_system_model', 'standard', 'abcmonitor', 'abcmonitor应用服务系统模式，总行 head 分行 branch 独立运行 standalone');
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_alarmengine_nfs_ip', '10.1.1.1', 'alarmengine', 'alarmengine nfs 挂载IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_alarmengine_replicas_number', '1', 'alarmengine', 'alarmengine deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_eventserver_replicas_number', '1', 'eventserver', 'eventserver deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_device-manager_mysql_ip', '10.1.1.1', 'device-manager', 'device-manager应用服务mysql的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_device-manager_replicas_number', '1', 'device-manager', 'device-manager deployment replicas 数量');
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_fengmap_replicas_number', '1', 'fengmap', 'fengmap deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_fengmap_nfs_ip', '10.1.1.1', 'fengmap', 'fengmap nfs 挂载IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_fileserver_filer_ip', '10.1.1.1', 'fileserver', 'fileserver weedfs的filer配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_fileserver_master_ip', '10.1.1.1', 'fileserver', 'fileserver weedfs的master配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_fileserver_replicas_number', '1', 'fileserver', 'fileserver deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_infra_mysql_ip', '10.1.1.1', 'infra', 'infra应用服务mysql的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_infra_system_model', 'standalone', 'infra', 'infra应用服务系统模式，总行 head 分行 branch 独立运行 standalone');
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_infra_replicas_number', '1', 'infra', 'infra deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_iot3rd_replicas_number', '1', 'iot3rd', 'iot3rd deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_iotserver_mysql_ip', '10.1.1.1', 'iotserver', 'iotserver' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_iotserver_weedfs_internal_ip', '10.1.1.1', 'iotserver', 'abcmonitor应用服务调取深瞳云的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_iotserver_weedfs_external_ip', '10.1.1.1', 'iotserver', 'abcmonitor应用服务调取深瞳云的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_iotserver_callback_url_ip', '10.1.1.1', 'iotserver', 'abcmonitor应用服务调取深瞳云的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_iotserver_iot_external_ip', '10.1.1.1', 'iotserver', 'abcmonitor应用服务调取深瞳云的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_iotserver_replicas_number', '1', 'iotserver', 'abcmonitor应用服务调取深瞳云的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_mediastream_ip', '10.1.1.1', 'mediastream', 'mediastream服务 nginx2的IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_mediastream_max_task', '10', 'mediastream', 'mediastream 视频路数' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_mediastream_replicas_number', '1', 'mediastream', 'mediastream deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_reportcenter_mysql_ip', '10.1.1.1', 'reportcenter', 'reportcenter应用服务mysql的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_reportcenter_visit_url_ip', '10.1.1.1', 'reportcenter', 'reportcenter visit url 一般为nginx2的IP');
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_reportcenter_replicas_number', '1', 'reportcenter', 'reportcenter deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_tts_replicas_number', '1', 'tts', 'tts deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_web_replicas_number', '1', 'web', 'web deployment replicas 数量' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_web_nfs_ip', '10.1.1.1', 'web', 'web nfs 挂载IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_flyway_mysql_ip', '10.1.1.1', 'flyway', 'flyway应用服务mysql的配置IP' );
INSERT INTO config ( name, value, app, comment ) VALUES ( 'tiantong_flyway_replicas_number', '1', 'flyway', 'deployment replicas 数量' );
