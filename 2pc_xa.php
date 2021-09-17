<?PHP

require "db.php";

$dbtest1 = new mysqli("127.0.0.1","root","root","test")or die("dbtest1 连接失败");
$dbtest2 = new mysqli("127.0.0.1","root","root","test_2")or die("dbtest2 连接失败");


//$dbtest1 = new db("127.0.0.1","root","root","test");
//$dbtest2 = new db("127.0.0.1","root","root","test_2");

//为XA事务指定一个id，xid 必须是一个唯一值。
$xid = uniqid("");

//两个库指定同一个事务id，表明这两个库的操作处于同一事务中
$dbtest1->query("XA START '$xid'");//准备事务1
$dbtest2->query("XA START '$xid'");//准备事务2


try {

    //$dbtest1

    $return = $dbtest1->query("UPDATE user SET age = 443  WHERE id = 1") ;
    echo "xa1:"; print_r($return);

    if(!in_array($return,['0','1'])) {
        throw new Exception("库1执行sql操作失败！");
    }

    //$dbtest2
    $return = $dbtest2->query("UPDATE user SET age = 443  WHERE id = 1") ;
    echo "xa2:"; print_r($return);
    if(!in_array($return,['0','1'])) {
        throw new Exception("库2执行sql操作失败！");
    }


    //阶段1：$dbtest1提交准备就绪
    $dbtest1->query("XA END '$xid'");

    $dbtest1->query("XA PREPARE '$xid'");

    //阶段1：$dbtest2提交准备就绪
    $dbtest2->query("XA END '$xid'");

    $dbtest2->query("XA PREPARE '$xid'");


    //阶段2：提交两个库
    $dbtest1->query("XA COMMIT '$xid'");

    $dbtest2->query("XA COMMIT '$xid'");

}

catch (Exception $e) {

    //阶段2：回滚

    $dbtest1->query("XA ROLLBACK '$xid'");
    /*
    上面这行代码是2pc中的xa事务，
    update set a = a+1
    如果是TCC，那么上面这行代码就变了，变成调用一个php接口，这个接口的作用就是把之前的操作给取消
    update set a = a-1*/
    $dbtest2->query("XA ROLLBACK '$xid'");


    die("Exception:".$e->getMessage());

}

echo "执行完毕";exit;

?>