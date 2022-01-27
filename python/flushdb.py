#!/usr/bin/python
import pymysql
import os
import redis

serverdb_ip     = "192.168.70.24"
serverdb_user   = "root"
serverdb_pass   = "admin"
serverdb_name   = "plat_db_test"

# 打开数据库连接
db = pymysql.connect(host=serverdb_ip,user=serverdb_user,password=serverdb_pass,database=serverdb_name)
# 使用 cursor() 方法创建一个游标对象 cursor
cursor = db.cursor()
sql = "SELECT MysqlIp,MysqlPassword,MysqlName,logDbName,ServerId,GobalRedisIp,GobalRedisPassword FROM ServerRegion"
# 使用 execute()  方法执行 SQL 查询
cursor.execute(sql)
# 使用 fetchall() 方法获取
alldata = cursor.fetchall()
# 关闭数据库连接
db.close()
#print(alldata)

if alldata:
    for data in alldata:
        print(str(data[4]) + " : 执行mysql==========>start")
        re_db = os.system("mysql -h" + data[0] + " -uroot -p" + data[1] + " -D" + data[2] +"<./clear_db.sql")
        print(re_db)
        re_log = os.system("mysql -h" + data[0] + " -uroot -p" + data[1] + " -D" + data[3] +"<./clear_log.sql")
        print(re_log)
        print(str(data[4]) + " : 执行mysql==========>end")
        
        print(str(data[4]) + " : 执行redis==========>start")
        re_con_redis = redis.Redis(host=data[5],port=6379,password=data[6],decode_responses=True,db=0)
        #re_con_redis = redis.Redis(host="127.0.0.1",port=6379,password="",decode_responses=True,db=0)
        re_redis     = re_con_redis.flushdb()
        print(re_redis)
        print(str(data[4]) + " : 执行redis==========>end\n\n")


#val = os.system("mysql -h" + serverdb_ip + " -u" + serverdb_user + " -p" + serverdb_pass + " -D" + serverdb_name +"<./clear.sql")
#print("mysql -h" + serverdb_ip + " -u" + serverdb_user + " -p" + serverdb_pass + " -D" + serverdb_name +"<./clear.sql")
#print(val)

