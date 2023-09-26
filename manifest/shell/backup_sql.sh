#!/bin/bash
#
# mysql docker container name
mysql="mysql-master"

# based today's date, create a folder in /backup/mysql
# and backup all databases to that folder
date=$(date +%Y-%m-%d)
back_dir="/backup/mysql/$date"

# create backup folder, if not exists, create it
if [ ! -d ${back_dir} ]; then
  mkdir -p ${back_dir}
fi

# 获取需要备份的数据库
docker exec ${mysql} mysql -u root -paDmin1b3@ -e "SHOW DATABASES;" | grep -Ev "Database|mysql|information_schema|performance_schema|sys" > ${back_dir}/database_list.txt

# Backup MySQL databases
while read -r database; do
    docker exec ${mysql} /usr/bin/mysqldump -u root --password=aDmin1b3@ "$database" > "${back_dir}/$database.sql"
done < ${back_dir}/database_list.txt

# compress backup folder
cd /backup/mysql
tar -zcvf ${date}.tar.gz ${date}

# delete backup folder
rm -rf ${back_dir}

# delete backup file 23 days ago
find /backup/mysql/ -mtime +23 -name "*.tar.gz" -exec rm -f {} \;


