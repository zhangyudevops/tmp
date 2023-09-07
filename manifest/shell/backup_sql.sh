#!/bin/bash
#
# incremental backup mysql data,sunday full backup, the other days of this week increment, and delete the backup files which are 21 days ago
# use mysqldump to backup mysql data

# define the backup path
backup_dir=/data/backup/mysql
# define the mysql data path
mysql_data_dir=/data/mysql/data
# define the mysql binlog path
mysql_binlog_dir=/data/mysql/binlog
# define the mysql user
mysql_user=root
# define the mysql password
mysql_password=123456
# define the mysql host
mysql_host=localhost
# define the mysql port
mysql_port=3306
# define the date format
date_format=`date +%Y%m%d`

# define the backup file name
backup_file_name=mysql-${date_format}.tar.gz

# define the backup log file name
backup_log_file_name=mysql_backup.log

# define the backup log file path
backup_log_file_path=${backup_dir}/${backup_log_file_name}

# define the backup file path
backup_file_path=${backup_dir}/${backup_file_name}

# define the mysql binlog file name
mysql_binlog_file_name=mysql-bin.index

# define the mysql binlog file path
mysql_binlog_file_path=${mysql_binlog_dir}/${mysql_binlog_file_name}

# define the mysql binlog backup file name
mysql_binlog_backup_file_name=mysql-bin-${date_format}.tar.gz

# define the mysql binlog backup file path
mysql_binlog_backup_file_path=${backup_dir}/${mysql_binlog_backup_file_name}

# define the mysql binlog backup log file name
mysql_binlog_backup_log_file_name=mysql_binlog_backup.log

