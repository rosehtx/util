<?php
/**
 * Created by PhpStorm.
 * User: htx
 * Date: 2021/3/17
 * Time: 23:27
 */

require "vendor/autoload.php";
date_default_timezone_set('PRC');
use Swoole\Coroutine;

use PhpAmqpLib\Connection\AMQPStreamConnection;
use PhpAmqpLib\Exchange\AMQPExchangeType;
use PhpAmqpLib\Message\AMQPMessage;

//====================amq
//$connection = new AMQPStreamConnection('192.168.182.129', 5672, 'guest', 'guest', '/');
//$channel = $connection->channel();
//
////$queue    = 'testConfirmQueue';
////$exchange = 'testConfirmExchange';
//$queue    = 'qos_queue';
//$exchange = 'router';
//$channel->queue_declare($queue, false, true, false, false);
////$channel->exchange_declare($exchange, AMQPExchangeType::FANOUT, false, true, false);
//$channel->exchange_declare($exchange, AMQPExchangeType::DIRECT, false, true, false);
//$channel->queue_bind($queue, $exchange);
//
//function process_message($message)
//{
//    echo ("\n--------\n");
//    echo ($message->body);
//    $data = $message->body;
////    echo ("\n--------\n");
////    go(function ()use($data) {
////        $res = file_put_contents("/home/wwwroot/default/test/1.log",$data.PHP_EOL,FILE_APPEND);
////    });
//    $message->ack(true);
//    // Send a message with the string "quit" to cancel the consumer.
//    if ($message->body === 'quit') {
//        $message->getChannel()->basic_cancel($message->getConsumerTag());
//    }
//}
//$consumerTag = 'consumer';
////go(function ()use($channel,$queue,$consumerTag) {
////    Coroutine::sleep(1);
//    $channel->basic_consume($queue, $consumerTag, false, false, false, false, 'process_message');
////});
//function shutdown($channel, $connection)
//{
//    $channel->close();
//    $connection->close();
//}
//register_shutdown_function('shutdown', $channel, $connection);
//// Loop as long as the channel has callbacks registered
//while ($channel->is_consuming()) {
//    $channel->wait();
//}

//====================kafka
//$config = \Kafka\ProducerConfig::getInstance();
//$config->setMetadataRefreshIntervalMs(10000);
//$config->setMetadataBrokerList('192.168.182.129:9093,192.168.182.129:9094');
//$config->setBrokerVersion('0.9.0.1');
//$config->setRequiredAck(1);
//$config->setIsAsyn(false);
//$config->setProduceInterval(500);
//$producer = new \Kafka\Producer();
//
//for($i = 0; $i < 1; $i++) {
//    $rr = $producer->send([
//        [
//            'topic' => 'test',
//            'value' => 'this is test',
//            'key' => '',
//        ],
//    ]);
//    var_dump($rr);
//}
//$config = \Kafka\ConsumerConfig::getInstance();
//$config->setMetadataRefreshIntervalMs(10000);
//$config->setMetadataBrokerList('192.168.182.129:9093,192.168.182.129:9094');
//$config->setGroupId('test1');
//$config->setBrokerVersion('0.9.0.1');
//$config->setTopics(['test']);
////$config->setOffsetReset('earliest');
//$consumer = new \Kafka\Consumer();
//$consumer->start(function($topic, $part, $message) {
//    var_dump($message);
//});

//====================redis集群
//$redis_clus = new RedisCluster(null,[
//    '192.168.182.129:6479',
//    '192.168.182.129:6480',
//    '192.168.182.129:6481',
//    '192.168.182.129:6482',
//    '192.168.182.129:6483',
//    '192.168.182.129:6484']);
//$u = $redis_clus->get('u');
//var_dump($u);

//====================mongodb操作
//$manager = new MongoDB\Driver\Manager("mongodb://192.168.182.129:27017");

// 插入数据
//$bulk = new MongoDB\Driver\BulkWrite;
//$bulk->insert(['x' => 1, 'name'=>'xiaohong', 'age' => '11']);
//$bulk->insert(['x' => 2, 'name'=>'小明', 'age' => '22']);

//$bulk->update(
//    ['x' => 2],
//    ['$set' => ['class' => ['语文','英语']]],
//    ['multi' => false, 'upsert' => false]
//);

//$bulk->update(
//    ['x' => 2],
//    ['$rename' => ['class' => 'object']]
//);
//
//$manager->executeBulkWrite('school.class', $bulk);

//====================elasticsearch
//$ch = curl_init();
// curl_setopt($ch,CURLOPT_URL,"http://192.168.182.127:9210/_bulk?pretty");
// curl_setopt($ch,CURLOPT_RETURNTRANSFER,1);
// curl_setopt($ch,CURLOPT_POST,1);
// curl_setopt($ch, CURLOPT_BINARYTRANSFER, 1);
// curl_setopt($ch,CURLOPT_POSTFIELDS,"@/home/wwwroot/default/test/test.json");
// curl_setopt($ch,CURLOPT_HTTPHEADER,array('Content-Type:application/json'));
//
// $result = curl_exec($ch);
// var_dump($result);

