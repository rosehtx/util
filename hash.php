<?php
//一致性HASH测试
class Consistance
{
    protected $num   = 10;   //设定每一个服务器的节点数，数量越多，宕机时服务器负载就会分布得越平均，但也增大数据查找消耗。
//    protected $nodes = [];   //当前服务器组的结点列表。
    public    $nodes = [];   //当前服务器组的结点列表。

    //计算一个数据的哈希值，用以确定位置
    public function make_hash($data)
    {
        // %u - 不包含正负号的十进制数（大于等于 0）
        // crc32 以整数值返回字符串的 32 位循环冗余校验码多项式
        return sprintf('%u', crc32($data));
    }

    //遍历当前服务器组的节点列表，确定需要存储/查找的服务器
    public function get_loc($data)
    {
        $loc = $this->make_hash($data);
        foreach ($this->nodes as $key => $val) {
            if ($loc <= $key) {
                return $val;
            }
        }
    }

    //添加一个服务器，将其结点添加到服务器组的节点列表内。
    public function add_host($host)
    {
        for ($i = 0; $i < $this->num; $i++) {
            //可以自己定义生成键值的方式，保证唯一
            $key               = $this->make_hash($host . '_' . $i);
            //节点只想对应服务器
            $this->nodes[$key] = $host;
        }
        ksort($this->nodes);        //对结点排序，这样便于查找。
    }

    //删除一个服务器，并将其对应节点从服务器组的节点列表内移除。
    public function remove_host($host)
    {
        for ($i = 0; $i < $this->num; $i++) {
            $key = $this->make_hash($host . '_' . $i);
            unset($this->nodes[$key]);
        }
    }
}

//代码测试
$Consistance = new Consistance();
$Consistance->add_host('servertest1');
$Consistance->add_host('servertest2');
$Consistance->add_host('servertest3');
//打印节点列表
print_r($Consistance->nodes);
echo '<hr>';
//echo对应服务器
echo $Consistance->get_loc('test_key1').'<br>'; //落在test3
echo $Consistance->get_loc('test_key2').'<br>'; //落在test1
echo '<hr>';
////去掉一个有落点服务器继续
//$Consistance->remove_host('servertest1');
//print_r($Consistance->nodes);
//echo '<hr>';
//echo $Consistance->get_loc('test_key1').'<br>'; //落在test3
//echo $Consistance->get_loc('test_key2').'<br>'; //落在test2

//去掉一个没有落点服务器继续
$Consistance->remove_host('servertest2');
print_r($Consistance->nodes);
echo '<hr>';
echo $Consistance->get_loc('test_key1').'<br>'; //落在test3
echo $Consistance->get_loc('test_key2').'<br>'; //落在test1


