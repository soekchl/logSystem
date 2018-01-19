# logSystem

# 调用模块添加
	socket - tcp
	传输用 pb
	创建 net 模块 专门用于传输log
	
这块的和 我 原先写的 mylog 要结合起来
	把MyLog弄到这里 合并
	
	

	1、增加分表机制



# docker 数据库配置

--- 

创建数据库(utf8)：
	CREATE DATABASE `logSystem` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci
	
创建帐号 admin-123456
	INSERT INTO `user` VALUES ('1', 'admin', 'c3584069f326b7e23ce4b6d45eaa9857', 'B8BvzRAKem', '', '', '2018-01-18 20:24:21', '[', '0', '2017-12-16 11:15:24');
	
	
	

创建docker 主数据库(3307)
	docker run --name master-mysql -p 3307:3306  -e MYSQL_ROOT_PASSWORD=admin -d mysql:5.7.20
	
创建docker 副数据库(3308)
	docker run --name slave-mysql -p 3308:3306 -e MYSQL_ROOT_PASSWORD=admin -d mysql:5.7.20

---

# Master配置文件
[mysqld]
## 设置server_id，一般设置为IP，同一局域网内注意要唯一
server_id=100  
## 复制过滤：也就是指定哪个数据库不用同步（mysql库一般不同步）
binlog-ignore-db=mysql  
## 开启二进制日志功能，可以随便取，最好有含义（关键就是这里了）
log-bin=edu-mysql-bin  
## 为每个session 分配的内存，在事务过程中用来存储二进制日志的缓存
binlog_cache_size=1M  
## 主从复制的格式（mixed,statement,row，默认格式是statement）
binlog_format=mixed  
## 二进制日志自动删除/过期的天数。默认值为0，表示不自动删除。
expire_logs_days=7  
## 跳过主从复制中遇到的所有错误或指定类型的错误，避免slave端复制中断。
## 如：1062错误是指一些主键重复，1032错误是因为主从数据库数据不一致
slave_skip_errors=1062
	
---

复制配置文件到 docker
	docker cp my.cnf master-mysql:/etc/mysql/conf.d/mysql.cnf
	
重启
	service mysql restart
	
创建同步用户
	CREATE USER 'slave'@'%' IDENTIFIED BY '123456';
	GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'slave'@'%';


---

# slave 配置
[mysqld]
## 设置server_id，一般设置为IP,注意要唯一
server_id=4  
## 复制过滤：也就是指定哪个数据库不用同步（mysql库一般不同步）
binlog-ignore-db=mysql  
## 开启二进制日志功能，以备Slave作为其它Slave的Master时使用
log-bin=edu-mysql-slave1-bin  
## 为每个session 分配的内存，在事务过程中用来存储二进制日志的缓存
binlog_cache_size=1M  
## 主从复制的格式（mixed,statement,row，默认格式是statement）
binlog_format=mixed  
## 二进制日志自动删除/过期的天数。默认值为0，表示不自动删除。
expire_logs_days=7  
## 跳过主从复制中遇到的所有错误或指定类型的错误，避免slave端复制中断。
## 如：1062错误是指一些主键重复，1032错误是因为主从数据库数据不一致
slave_skip_errors=1062  
## relay_log配置中继日志
relay_log=edu-mysql-relay-bin  
## log_slave_updates表示slave将复制事件写进自己的二进制日志
log_slave_updates=1  
## 防止改变数据(除了特殊的线程)
read_only=1

---


复制配置文件 和 从起


master 中运行和查看  开始备份的点
	show master status;


slave的mysql中运行
	change master to master_host='172.17.0.2', master_user='slave', master_password='123456', master_port=3306, master_log_file='edu-mysql-bin.000001', master_log_pos=617, master_connect_retry=30;
				这里从上面获得注意


查看 Slavel 状态
	show slave status \G;

	正常状态	
	Slave_IO_Running: Yes
	Slave_SQL_Running: Yes

启动 Slavel
	start slave;










