V1.5
1. done [new Feature] 支持模糊查询；
2. done [new Feature] 支持导入导出功能；
3. done [fix] 操作记录返回按照时间排序，最新的在前面

V1.4
1. [new Feature] 各种前端接口支持；

V1.3

1. [New Feature] 增加数据库词库变更，定时轮询对比加载功能
2. [bugfix]  默认表明为green_shield_words_infos，导致找不到表；修改为手动指定表名；
3. [todo] 增加./green_shield -v 显示版本号和使用帮助的功能；

V1.2

1. [fixed]预发和线上环境连接不上数据库; tcp(xxxxx:3306) 连接地址要加上tcp()为啥，不清楚；测试环境也不需要；
2. [fixed]日志打印格式不符合预期；(logrus的text formatter的实现有两种，带color的和不带color的，格式不一样，如果输出到stdout, 那么带color的格式）

