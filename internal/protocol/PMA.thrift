namespace go PMA

/* thrift -r --gen go PMA.thrift */


/* 防误代理内容实体 */
struct PMAMsg {
    /* 请求头信息*/
    1: required map<string,string> head;

    /* 源账号 */
    2: required string src;

    /* 目标账号 */
    3: required list<string> targets;

    /* 实际业务数据 */
    4: required string content;

}


/* 防误代理服务 */
service PMAService {
    /*
    * 请求消息方法，可以是子站发防误代理，也可以是防误代理发子站
    * 参数：
    * pmaMsg：内容实体
    */
    oneway void requestFunc(1:PMAMsg pmaMsg);
}
