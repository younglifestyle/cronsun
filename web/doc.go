package web

/**
* @api GET /v1/version 版本号
* @apiGroup worker

* @apiSuccess 200 OK
* @apiExample string
* v0.0.5 (build go1.15)
*/

/**
* @api GET /v1/jobs 获取任务列表
* @apiGroup worker
* @apiQuery group    string   限制特定的任务组
* @apiQuery node     string   限制特定的节点

* @apiSuccess 200 OK
* @apiExample json
[
    {
        "id": "12dc3e3f",
        "name": "test",
        "group": "default",
        "cmd": "echo \"12414\"",
        "run_param": "",
        "rules": [
            {
                "id": "da7fdade10ff00bbeb633024bef7ab0aa7be9419448f2fd3",
                "timer": "0/1 * * * * ?",
                "gids": null,
                "nids": [
                    "a4e730ee-4834-4981-a434-4620bcd46d834"
                ],
                "exclude_nids": null
            }
        ],
        "pause": false,
        "timeout": 0,
        "parallels": 1,
        "retry": 0,
        "interval": 0,
        "kind": 0,
        "avg_time": 0,
        "fail_notify": false,
        "to": [],
        "log_expiration": 0,
        "cmd_type": "SHELL",
        "update_time": 1598860512,
        "latestStatus": {
            "id": "5f4cac84ae6a0189255d7209",
            "jobId": "12dc3e3f",
            "jobGroup": "default",
            "user": "",
            "name": "test",
            "node": "a4e730ee-4834-4981-a434-4620bcd46d834",
            "hostname": "Z-CSH-2202B",
            "ip": "172.17.100.145",
            "success": true,
            "beginTime": "2020-09-01T17:02:04+08:00",
            "endTime": "2020-09-01T17:02:04.007+08:00",
            "refLogId": "5f4e0e0ca4c342434d129c0f"
        },
        "nextRunTime": "2020-09-01 17:02:05"
    }
]
*/

/**
* @api GET /v1/all/job/groups 获取所有任务分组
* @apiGroup worker

* @apiSuccess 200 OK
* @apiExample array
[
    "default"
]
*/

/**
* @api GET /v1/job/groups 获取所有任务分组
* @apiGroup worker

* @apiSuccess 200 OK
* @apiExample array
[
    "default"
]
*/

/**
* @api PUT /v1/job 创建或者更新任务信息
* @apiGroup worker

 * @apiRequest json
* @apiExample json
{
    "id": "7c20e504",    // 任务的内部ID
    "name": "test1",     // 任务名
    "group": "test",     // 命令所属的任务分组，利于查看及分类
    "cmd": "sleep 9s",   // 需要执行的Shell命令
    "cmd_type": "SHELL", // 执行的任务类型，可支持SHELL、Python
    "run_param": "",     // 给脚本执行的参数，多参数使用逗号分隔
   "rules": [
        {
            "id": "43dd2843",                  // 执行规则的内部ID
            "timer": "0/1 * * * * ? ",         // 任务的定时规则
            "gids": [],                        // 选择在哪些<节点组>运行
            "nids": [                          // 选择在哪些<节点>运行
                "a4e730ee-4834-4981-a434-4620bcd46d82",
                "a4e730ee-4834-4981-a434-4620bcd46d83"
            ],
            "exclude_nids": null               // 选择此规则<不在>哪些<节点>运行
        }
    ],
    "pause": false,                            // 暂停任务执行
    "timeout": 0,                              // 任务执行时间超时设置，大于0则生效
    "parallels": 0,                            // 任务是否可并行执行（例如单机状态，只要定时时间到，就开启任务，此时可能会出现并发执行任务）
    "retry": 0,                                // 任务失败重试次数
    "interval": 0,                             // 任务失败重试间隔
    "kind": 1,                                 // 任务类型，0: 普通任务；1: 单机任务，2：节点组内任务；普通任务即各可执行此任务的节点到达执行时间，均可执行；单机任务则是分布式执行，只能单机执行；
    "avg_time": 0,                             // 任务平均执行时间
    "fail_notify": false, 					   // 失败邮件通知
    "to": [],								   // 通知的地址
    "log_expiration": 0,					   // 任务日志是否清除(需要使能任务清理)
    "oldGroup": "default"					   // 之前的任务分组
}

* @apiSuccess 200 OK
* @apiExample null
""
*/

/**
* @api POST /v1/job/:group/:id?pause= 开始和暂停任务执行
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery id     string   任务ID
* @apiQuery pause     string   默认不填，则为false；填写则为true

* @apiSuccess 200 OK
* @apiExample 修改后的Job信息
{
    "id": "7c20e504",    // 任务的内部ID
    "name": "test1",     // 任务名
    "group": "test",     // 命令所属的任务分组，利于查看及分类
    "cmd": "sleep 9s",   // 需要执行的Shell命令
    "cmd_type": "SHELL",   // 执行的任务类型，可支持SHELL、Python
    "rules": [
        {
            "id": "43dd2843",                  // 执行规则的内部ID
            "timer": "0/1 * * * * ? ",         // 任务的定时规则
            "gids": [],                        // 选择在哪些<节点组>运行
            "nids": [                          // 选择在哪些<节点>运行
                "a4e730ee-4834-4981-a434-4620bcd46d82",
                "a4e730ee-4834-4981-a434-4620bcd46d83"
            ],
            "exclude_nids": null               // 选择此规则<不在>哪些<节点>运行
        }
    ],
    "pause": true,                            // 暂停任务执行
    "timeout": 0,                              // 任务执行时间超时设置，大于0则生效
    "parallels": 0,                            // 任务是否可并行执行（例如单机状态，只要定时时间到，就开启任务，此时可能会出现并发执行任务）
    "retry": 0,                                // 任务失败重试次数
    "interval": 0,                             // 任务失败重试间隔
    "kind": 1,                                 // 任务类型，0: 普通任务；1: 单机任务，2：节点组内任务；普通任务即各可执行此任务的节点到达执行时间，均可执行；单机任务则是分布式执行，只能单机执行；
    "avg_time": 0,                             // 任务平均执行时间
    "fail_notify": false, 					   // 失败邮件通知
    "to": [],								   // 通知的地址
    "log_expiration": 0,					   // 任务日志是否清除(需要使能任务清理)
    "oldGroup": "default"					   // 之前的任务分组
}
*/

/**
* @api POST /v1/job/:op 批量开始或暂停任务
* @apiGroup worker
* @apiQuery op  string   [start/stop]操作

 * @apiRequest json
* @apiExample array
["default/12dc3e3f","default/13dc3e3f"]

* @apiSuccess 200 OK
* @apiExample string
"2 of 2 updated."
*/

/**
* @api GET /v1/job/:group/:id 获取单个任务信息
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery id     string   任务ID

* @apiSuccess 200 OK
* @apiExample json
{
    "id": "7c20e504",    // 任务的内部ID
    "name": "test1",     // 任务名
    "group": "test",     // 命令所属的任务分组，利于查看及分类
    "cmd": "sleep 9s",   // 需要执行的Shell命令
    "cmd_type": "SHELL",   // 执行的任务类型，可支持SHELL、Python
    "rules": [
        {
            "id": "43dd2843",                  // 执行规则的内部ID
            "timer": "0/1 * * * * ? ",         // 任务的定时规则
            "gids": [],                        // 选择在哪些<节点组>运行
            "nids": [                          // 选择在哪些<节点>运行
                "a4e730ee-4834-4981-a434-4620bcd46d82",
                "a4e730ee-4834-4981-a434-4620bcd46d83"
            ],
            "exclude_nids": null               // 选择此规则<不在>哪些<节点>运行
        }
    ],
    "pause": true,                            // 暂停任务执行
    "timeout": 0,                              // 任务执行时间超时设置，大于0则生效
    "parallels": 0,                            // 任务是否可并行执行（例如单机状态，只要定时时间到，就开启任务，此时可能会出现并发执行任务）
    "retry": 0,                                // 任务失败重试次数
    "interval": 0,                             // 任务失败重试间隔
    "kind": 1,                                 // 任务类型，0: 普通任务；1: 单机任务，2：节点组内任务；普通任务即各可执行此任务的节点到达执行时间，均可执行；单机任务则是分布式执行，只能单机执行；
    "avg_time": 0,                             // 任务平均执行时间
    "fail_notify": false, 					   // 失败邮件通知
    "to": [],								   // 通知的地址
    "log_expiration": 0,					   // 任务日志是否清除(需要使能任务清理)
	"cmd_type": "SHELL",                       // 命令类型
    "update_time": 1598963384                  // 创建或更新时间
}
*/

/**
* @api DELETE /v1/job/:group/:id 删除单个任务
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery id     string   任务ID

* @apiSuccess 200 OK
* @apiExample json
“”
*/

/**
* @api GET /v1/job/:group/:id/nodes 获取执行该任务的节点
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery id     string   任务ID

* @apiSuccess 200 OK
* @apiExample json
[
    "a4e730ee-4834-4981-a434-4620bcd46d834"
]

*/

/**
* @api GET /v1/job/:group/:id/nodes 获取执行该任务的节点
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery id     string   任务ID

* @apiSuccess 200 OK
* @apiExample array
[
    "a4e730ee-4834-4981-a434-4620bcd46d834"
]

*/

/**
* @api GET /v1/job/:group/:id/nodes 获取执行该任务的节点
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery id     string   任务ID

* @apiSuccess 200 OK
* @apiExample array
[
    "a4e730ee-4834-4981-a434-4620bcd46d834"
]

*/

/**
* @api PUT /v1/job/:group/:id/execute?node= 立即执行任务（指定单个节点或是所有节点）
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery job     string   任务ID
* @apiQuery node     string   节点ID,为空则所有节点执行任务

* @apiSuccess 200 OK
* @apiExample string
“”
*/

/**
* @api GET /v1/job-executing 获取执行中的任务信息
* @apiGroup worker
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery job     string   任务ID
* @apiQuery node     string   节点ID,为空则所有节点执行任务

* @apiSuccess 200 OK
* @apiExample json
[
    {
        "id": "25405",                      // job PID
        "jobId": "12dc3e3f",                // job id
        "group": "default",                 // 任务分组信息
        "nodeId": "a4e730ee-4834-4981-a434-4620bcd46d834",   // 节点ID
        "time": "2020-09-01T21:16:43.0003405+08:00",         // 执行时间
        "killed": false,                                     // 是否强杀
        "jobName": "test"                                    // 任务名
    }
]

*/

/**
* @api DELETE /v1/job-executing 强杀任务
* @apiGroup worker
* @apiQuery node  string   节点ID
* @apiQuery group  string   任务分组信息，默认填入[default]
* @apiQuery job     string   任务ID
* @apiQuery pid     string   任务执行的PID

* @apiSuccess 200 OK
* @apiExample string
"Killing process"

*/

/**
* @api GET /v1/logs 获取执行中的任务信息
* @apiGroup worker
* @apiQuery hostnames  []string  主机的名字集合
* @apiQuery ips     []string   主机的IP地址集
* @apiQuery names     []string   任务名集合
* @apiQuery ids     []string   任务ID集合
* @apiQuery begin     string   时间区间，2020-08-30
* @apiQuery end     string   时间区间，2020-08-31
* @apiQuery page     string   当前页
* @apiQuery pageSize     int   每页显示数量，默认显示50
* @apiQuery failedOnly     bool   节点ID,为空则所有节点执行任务

* @apiSuccess 200 OK
* @apiParam total int     总数量
* @apiParam list json     log信息
* @apiExample json
{
    "total": 270,
    "list": [
        {
            "id": "5f4ca949a4c3422a1267d6bf",
            "jobId": "13dc3e3f",
            "jobGroup": "default",
            "user": "",
            "name": "test",
            "node": "a4e730ee-4834-4981-a434-4620bcd46d834",
            "hostname": "Z-CSH-2202B",
            "ip": "172.17.100.145",
            "success": false,
            "beginTime": "2020-08-31T15:39:53+08:00",
            "endTime": "2020-08-31T15:39:53.007+08:00"
        },
        {
            "id": "5f4ca948a4c3422a1267d6be",
            "jobId": "13dc3e3f",
            "jobGroup": "default",
            "user": "",
            "name": "test",
            "node": "a4e730ee-4834-4981-a434-4620bcd46d834",
            "hostname": "Z-CSH-2202B",
            "ip": "172.17.100.145",
            "success": false,
            "beginTime": "2020-08-31T15:39:52+08:00",
            "endTime": "2020-08-31T15:39:52.006+08:00"
        }
]

*/

/**
* @api GET /v1/log/:id 获取指定任务的最新一条Log信息
* @apiGroup worker
* @apiQuery id     string   任务ID

* @apiSuccess 200 OK
* @apiParam total int     总数量
* @apiParam list json     log信息
* @apiExample json
{
    "id": "5f4e4ccea4c3425fa589c0c4",
    "jobId": "13dc3e3f",
    "jobGroup": "default",
    "user": "",
    "name": "test",
    "node": "a4e730ee-4834-4981-a434-4620bcd46d834",
    "hostname": "Z-CSH-2202B",
    "ip": "172.17.100.145",
    "command": "#!/bin/bash\necho \"xxl-job: hello shell\"\n\necho \"脚本位置：$0\"\necho \"任务参数：$1\"\necho \"分片序号 = $2\"\necho \"分片总数 = $3\"\n \necho \"Good bye!\"\nexit 0",
    "output": "xxl-job: hello shell\n脚本位置：/data/cron-job/filesource/13dc3e3f_1598860120.sh\n任务参数：\n分片序号 = \n分片总数 = \nGood bye!\n",
    "success": true,
    "beginTime": "2020-09-01T21:29:50+08:00",
    "endTime": "2020-09-01T21:29:50.011+08:00"
}

*/

/**
* @api GET /v1/nodes 获取所有节点信息
* @apiGroup worker

* @apiSuccess 200 OK
* @apiExample json
[
    {
        "id": "a4e730ee-4834-4981-a434-4620bcd46d834",
        "pid": "24485",
        "ip": "172.17.100.145",
        "hostname": "Z-CSH-2202B",
        "version": "v0.3.5 (build go1.15)",
        "up": "2020-09-01T21:02:32.315+08:00",
        "down": "0001-01-01T00:00:00Z",
        "alived": true,
        "connected": true
    }
]
*/

/**
* @api DELETE /v1/node?id= 删除故障或者离线节点信息
* @apiGroup worker
* @apiQuery id     string   节点ID

* @apiSuccess 200 OK
* @apiExample json
""
*/

/**
* @api DELETE /v1/node/groups 获取节点所有组信息
* @apiGroup worker

* @apiSuccess 200 OK
* @apiExample json
[
    {
        "id": "7b6cf77f02b26949ff61f91d5fdb41ee524f13cb28c59030",    // Group ID
        "name": "test",                                              // Group组名
        "nids": [
            "a4e730ee-4834-4981-a434-4620bcd46d834"    // 节点ID
        ]
    }
]

*/

/**
* @api GET /v1/node/group?id= 通过Group ID获取节点组信息
* @apiGroup worker
* @apiQuery id     string   节点ID

* @apiSuccess 200 OK
* @apiExample json
{
    "id": "7b6cf77f02b26949ff61f91d5fdb41ee524f13cb28c59030",
    "name": "test",
    "nids": [
        "a4e730ee-4834-4981-a434-4620bcd46d834"
    ]
}
*/

/**
* @api PUT /v1/node/group 创建或更新一个节点组
* @apiGroup worker
* @apiQuery id     string   节点ID

 * @apiRequest json
 * @apiParam name    string	   节点组名
 * @apiParam nids     []string  节点ID集
* @apiExample json
{
    "name": "weq1",
    "nids": [
        "a4e730ee-4834-4981-a434-4620bcd46d834"
        "a4e730ee-4834-4981-a434-4620bcd46d124"
    ]
}

* @apiSuccess 200 OK
* @apiExample json
""
*/

/**
* @api DELETE /v1/node/group?id= 删除一个节点组
* @apiGroup worker
* @apiQuery id     string   节点ID

* @apiSuccess 200 OK
* @apiExample json
""
*/

/**
* @api PUT /v1/info/overview 获取最近worker节点运行情况
* @apiGroup worker

* @apiSuccess 200 OK
* @apiExample json
{
    "totalJobs": 2,
    "jobExecuted": {
        "total": 237700,
        "successed": 213942,
        "failed": 23758,
        "date": ""
    },
    "jobExecutedDaily": [
        {
            "total": 44284,
            "successed": 44284,
            "failed": 0,
            "date": "2020-08-26"
        },
        {
            "total": 0,
            "successed": 0,
            "failed": 0,
            "date": "2020-08-27"
        },
        {
            "total": 9226,
            "successed": 2030,
            "failed": 7196,
            "date": "2020-08-28"
        },
        {
            "total": 0,
            "successed": 0,
            "failed": 0,
            "date": "2020-08-29"
        },
        {
            "total": 0,
            "successed": 0,
            "failed": 0,
            "date": "2020-08-30"
        },
        {
            "total": 17614,
            "successed": 4140,
            "failed": 13474,
            "date": "2020-08-31"
        },
        {
            "total": 12641,
            "successed": 9856,
            "failed": 2785,
            "date": "2020-09-01"
        }
    ]
}
*/
