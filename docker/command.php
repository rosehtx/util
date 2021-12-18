<?php
//swoft
//docker run -itd --restart=always -p 18306:18306 -p 18317:18307 -v /home/wwwroot/default/docker/swoft_goods:/var/www/swoft --network=swoft --ip 182.158.0.10 --name swoft_10 swoft:v1 sh
//docker exec -it swoft_10 php swoft/bin/swoft rpc:start ext_init=-i:192.168.182.127?-t:tcp?-p:18317?-n:goods
//docker exec -it swoft_20 php swoft/bin/swoft rpc:start ext_init=-i:192.168.182.127?-t:tcp?-p:18318?-n:goods
//docker run -itd --restart=always -p 18336:18306 -p 18347:18307 -v /home/wwwroot/default/docker/swoft_order:/var/www/swoft --network=swoft --ip 182.158.0.15 --name swoft_order_15 swoft:v1 sh
//docker exec -it swoft_order_15 php swoft/bin/swoft rpc:start ext_init=-i:192.168.182.127?-t:tcp?-p:18347?-n:order
//docker exec -it swoft_order_16 php swoft/bin/swoft rpc:start ext_init=-i:192.168.182.127?-t:tcp?-p:18348?-n:order


//consul
//docker run -itd --restart=always -p 8510:8500 --network=swoft --ip 182.158.0.110 --name consul_8510_110 consul1.4 ./consul agent -server -bootstrap-expect 1 -data-dir /tmp/consul -node=ali -bind=182.158.0.110 -ui -client=0.0.0.0

//rabbitmq
//docker run -itd --restart=always -p 15672:15672 -p 5672:5672 --hostname mq1 --network=rabbitmq --ip 172.148.0.10 --name rabbitmq_10 rabbitmq:3.7-management-alpine

//haproxy(这里要配置映射的路径，因为没有配置导致无法进入到docker)
//docker run -itd --restart=always -p 8100:8100 -p 5600:5600 -v /home/wwwroot/default/docker/haproxy/haproxy_10:/haproxy --privileged --network=haproxy --ip 172.147.0.10 --name haproxy_10 haproxy:1.7.13-alpine haproxy -f /haproxy/haproxy.cfg

//mongodb-分片集容器
//docker run -itd --restart=always --expose=27017 --expose=27018 --expose=27019 -m 1G --privileged --network=mongodb --ip 172.146.0.10 --name mongo_shard10 mongo:4.0 mongod --shardsvr --directoryperdb --replSet shard1
//mongodb-配置节点
//docker run -itd --restart=always --expose=27017 --expose=27018 --expose=27019 -m 1G --privileged --network=mongodb --ip 172.146.0.20 --name mongo_config20 mongo:4.0 mongod --configsvr --directoryperdb --replSet mongo-config --smallfiles
//mongodb-router节点
//docker run -itd --restart=always --expose=27017 --expose=27018 --expose=27019 -m 1G --privileged --network=mongodb --ip 172.146.0.30 --name mongo_router30 mongo:4.0 mongos --configdb mongo-config/172.146.0.20:27019,172.146.0.21:27019 --bind_ip 0.0.0.0 --port 27017
//mongo，config服务配置
//docker exec -it mongo_config20 bash -c "echo 'rs.initiate({_id: \"mongo-config\",configsvr: true, members: [{ _id : 0, host : \"172.146.0.20:27019\" },{ _id : 1, host : \"172.146.0.21:27019\" }]})' | mongo --port 27019"
//mongo 分片服配置
//docker exec -it mongo_shard10 bash -c "echo 'rs.initiate({_id: \"shard1\",members: [{ _id : 10, host : \"172.146.0.10:27018\" },{ _id : 11, host : \"172.146.0.11:27018\" },{ _id : 12, host : \"172.146.0.12:27018\", arbiterOnly: true }]})' | mongo --port 27018"
//router服务添加
//docker exec -it mongo_router30 bash -c "echo 'sh.addShard(\"shard1/172.146.0.10:27018,172.146.0.11:27018,172.146.0.12:27018\")' | mongo"

//redis
//docker run -itd --restart=always --network=host --name redis_cluster_200 rosehtx/nginx_php_redis_swoole:v3 sh /data/up.sh
//redis-cli --cluster create 192.168.182.129:6479 192.168.182.129:6480 192.168.182.129:6481 192.168.182.129:6482 192.168.182.129:6483 192.168.182.129:6484  --cluster-replicas 1
//mkdir /redis && mkdir /redis/data && mkdir /redis/log && touch /redis/log/redis.log

//mysql
//docker run -itd --restart=always --expose=22 -p 3316:3306 --network=dockernet1 --ip 192.168.0.36 --name mysql_master_36 rosehtx/nginx_php_redis_swoole:v4 sh /data/up.sh
//change master to master_host='192.168.182.127',
//master_port=3316,
//master_user='root',
//master_password='root',
//master_log_file='mysql-bin.000002',
//master_log_pos=0;

//elasticsearch(主)
//docker run -itd -p 9210:9200 -p 9310:9300 --network=dockernet1 --ip 192.168.0.40 -e ES_JAVA_OPTS="-Xms512m -Xmx512m" -v /home/wwwroot/default/docker/elasticsearch/master_10/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml \
//-v /home/wwwroot/default/docker/elasticsearch/master_10/data:/usr/share/elasticsearch/data \
//--privileged --name es_master_40 elasticsearch:7.6.1
//kibana管理elasticsearch
//docker run -itd -p 5601:5601 --network=dockernet1 --ip 192.168.0.45 -e ELASTICSEARCH_HOSTS="http://192.168.0.40:9200" -e SERVER_PORT="5601" -e SERVER_HOST="0.0.0.0" --name kibana kibana:7.6.1

//zookeeper
//docker run -itd -p 2182:2181 --network=dockernet1 --ip 192.168.0.82 -e ZOO_MY_ID="1" -e ZOO_SERVERS="server.1=192.168.0.82:2888:3888 server.2=192.168.0.83:2888:3888" --name zoo1 zookeeper:3.4.14 bash
//docker run -itd -p 2183:2181 --network=dockernet1 --ip 192.168.0.83 -e ZOO_MY_ID="2" -e ZOO_SERVERS="server.1=192.168.0.82:2888:3888 server.2=192.168.0.83:2888:3888" --name zoo2 zookeeper:3.4.14 bash
//kafka（包括可视化管理容器的安装）
//kafka1
//docker run -itd -p 9093:9092 --network=dockernet1 --ip 192.168.0.93 -h="kafka1" --link=192.168.0.82 --link=192.168.0.83 --name kafka1 wurstmeister/kafka:2.12-2.4.0 bash
//kafka2
//docker run -itd -p 9094:9092 --network=dockernet1 --ip 192.168.0.94 -h="kafka2" --link=192.168.0.82 --link=192.168.0.83 --name kafka2 wurstmeister/kafka:2.12-2.4.0 bash
//kafka-manager,不同服务器可以去掉--link
//docker run -itd -p 9001:9000 --network=dockernet1 --ip 192.168.0.91 -e ZK_HOSTS="192.168.0.82:2181,192.168.0.83:2181" -e TZ="CST-8" \
//--link=192.168.0.93 --link=192.168.0.94 --link=192.168.0.82 --link=192.168.0.83 --name kafka-manager sheepkiller/kafka-manager:latest bash



