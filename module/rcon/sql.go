package rcon

var databaseSql = []string{
	`create table admin
(
    a_id        int auto_increment comment '管理员ID'
        primary key,
    ar_id       int          default 0   not null comment '角色ID',
    a_username  varchar(20)              not null comment '用户名',
    a_password  char(32)                 not null comment '用户密码',
    a_nickname  varchar(30)              not null comment '昵称',
    a_status    int          default 1   not null comment '账户状态,1->正常,4->删除',
    last_time   int unsigned default '0' not null comment '上次登录时间,[yc:time]',
    create_time int unsigned             not null comment '创建时间,[yc:time]'
);

create table admin_auth
(
    aa_id   int auto_increment
        primary key,
    aa_type int default 1 not null comment '权限类型,1->路由,2->菜单/按钮',
    aa_pid  int default 0 not null comment '权限的上级ID',
    aa_name varchar(30)   null comment '权限名',
    aa_path varchar(100)  not null comment '权限路径'
)
    comment '管理员权限组';

create table admin_log
(
    al_id       int auto_increment
        primary key,
    a_id        int          not null,
    al_content  varchar(200) not null comment '操作内容',
    al_type     int          not null comment '操作类型,1->登录,2->命令操作,3->配置文件,4->接口请求',
    create_time int          not null comment '操作时间,[ys:time]',
    al_tag      char(15)     not null comment '操作标签'
);

create index admin_log_a_id_index
    on admin_log (a_id);

create table admin_role
(
    ar_id    int auto_increment
        primary key,
    ar_name  varchar(20) not null comment '权限角色名',
    ar_group text        not null comment '权限组'
)
    comment '权限角色';

create table game_bill
(
    gb_id       int auto_increment
        primary key,
    gb_type     int unsigned  not null comment '记录类型,1->玩家击杀,2->对局结算,3->新对局,4->玩家进服,5->建队,6->上帝视角,7->玩家退服,8->签到,9->玩家伤害,10->清除数据',
    atk_u_id    int default 0 not null comment '攻击者uid',
    victim_u_id int default 0 not null comment '被攻击玩家uid',
    gb_msg      varchar(200)  not null comment '日志内容',
    gt_id       int           not null comment '对局标签',
    log_time    int           not null comment '日志时间,[ys:time]',
    create_time int           not null comment '记录时间,[ys:time]'
);

create index game_bill_atk_u_id_index
    on game_bill (atk_u_id);

create index game_bill_victim_u_id_index
    on game_bill (victim_u_id);

create table game_tag
(
    gt_id       int auto_increment
        primary key,
    gt_map      varchar(50)  not null comment '地图类名',
    gt_layer    varchar(50)  not null comment '地图具体类',
    log_time    int unsigned not null comment '日志记录时间,[ys:time]',
    create_time int unsigned not null comment '创建时间,[ys:time]'
)
    comment '对局标签';

create table server_auth
(
    sa_id    int auto_increment
        primary key,
    sa_name  char(20) not null comment '服务器配置的权限名',
    sa_value char(25) not null comment '权限的值'
);

create table server_role
(
    sr_id       int auto_increment
        primary key,
    sr_name     char(30)                not null comment '权限组的译名',
    sr_value    char(30)                not null comment '权限组的英文名',
    sr_auth     varchar(400) default '' not null comment '权限集合',
    create_time int                     not null comment '创建时间,[ys:time]'
);

create table system_setting
(
    ss_user_max    tinyint unsigned default '100' not null comment '最大用户数',
    ss_vip_reserve tinyint unsigned               not null comment '给VIP用户的预留位',
    ss_op_reserve  int              default 0     not null comment 'op的预留位,默认0,填入则保证至少有填入数量的op可进服'
)
    comment '系统设置';

create table user
(
    u_id         int auto_increment comment '用户ID'
        primary key,
    u_name       varchar(40)              not null comment '昵称',
    u_steam      varchar(32)              not null comment 'steamId',
    u_eos        varchar(32)              not null comment '服务器EosId',
    u_vip_level  tinyint      default 0   not null comment '用户的VIP等级',
    u_vip_expire int unsigned default '0' not null comment '用户会员过期时间,[ys:time]',
    u_black_info varchar(30)  default ''  not null comment '封禁原因',
    online_time  int unsigned default '0' not null comment '累计在线时长,[ys:time]',
    last_time    int unsigned default '0' not null comment '最后一次进服时间,[ys:time]',
    create_time  int unsigned             not null comment '创建时间,[ys:time]',
    black_time   int          default 0   not null comment '封禁时间,[ys:time]',
    sr_id        int          default 0   not null comment '权限组',
    sr_expire    int          default 0   not null comment '玩家服务器到期时间,[ys:time]'
)
    comment '用户表';

create index user_u_eos_index
    on user (u_eos);

create index user_u_steam_index
    on user (u_steam);

create table user_game_chess
(
    u_id       int                      not null comment '用户ID'
        primary key,
    ugc_kill   int unsigned default '0' not null comment '累计杀敌数',
    ugc_die    int unsigned default '0' not null comment '死亡数',
    ugc_rescue int unsigned default '0' not null comment '营救他人次数'
)
    comment '用户对局信息';

create table user_point
(
    u_id        int           not null
        primary key,
    u_points    int unsigned  not null comment '积分',
    points_time int default 0 not null comment '签到时间,[ys:time]'
)
    comment '用户点数信息';
`,
	`INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (1, '更换地图', 'changemap');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (2, '暂停游戏', 'pause');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (3, '作弊命令', 'cheat');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (4, '服务器密码', 'private');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (5, '忽略队伍平衡', 'balance');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (6, '超管聊天/广播', 'chat');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (7, '踢人', 'kick');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (8, '黑名单', 'ban');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (9, '服务器配置', 'config');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (10, '观察视角', 'cameraman');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (11, '防踢/ban', 'immune');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (12, '关闭服务器', 'manageserver');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (13, '预留位', 'reserve');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (14, '录制演示', 'demos');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (15, '回放演示', 'clientdemos');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (16, '无限更换队伍', 'teamchange');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (17, '跳边命令', 'forceteamchange');
INSERT INTO server_auth (sa_id, sa_name, sa_value) VALUES (18, '查看超管聊天/通知', 'canseeadminchat');
INSERT INTO admin (a_id, ar_id, a_username, a_password, a_nickname, a_status, last_time, create_time) VALUES (1, 1, 'admin', '123456', '默认管理员', 1, 0, 0);
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (1, 1, 0, '对局信息', 'chat');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (2, 1, 0, '对局数据', 'data');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (3, 1, 0, '管理员', 'admin');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (4, 1, 0, '玩家管理', 'online');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (5, 1, 0, '设置', 'setting');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (6, 2, 3, '创建账户', 'create');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (7, 2, 3, '禁用', 'disabled');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (8, 2, 3, '编辑', 'edit');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (9, 1, 0, '仪表盘', 'panel');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (10, 2, 9, '右键-解散队伍', 'dissolveSquad');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (11, 2, 9, '右键-移出队伍', 'removeSquad');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (12, 2, 9, '右键-踢人', 'kick');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (13, 2, 9, '右键-警告', 'warn');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (14, 1, 9, '右键-跳边', 'changeTeam');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (15, 1, 9, '右键-封禁', 'ban');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (16, 2, 1, '命令模块', 'command');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (17, 2, 1, '命令模块-广播', 'command_broadcast');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (18, 2, 1, '命令模块-下局自动配平', 'command_random_next');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (19, 2, 1, '命令模块-立即配平', 'command_random_now');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (20, 2, 1, '命令模块-跳边', 'command_changeTeam');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (23, 2, 1, '右键-踢人', 'click_kick');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (24, 2, 1, '右键-警告', 'click_warn');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (25, 2, 1, '右键-跳边', 'click_changeTeam');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (26, 2, 1, '右键-封禁', 'click_ban');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (27, 1, 0, '命令输入框', 'command_data');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (28, 2, 27, '踢出玩家', 'AdminKick');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (29, 2, 27, '广播', 'AdminBroadcast');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (30, 2, 27, '立即结束当前比赛', 'AdminEndMatch');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (31, 2, 27, '改变地图', 'AdminChangeLayer');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (32, 2, 27, '改变下一张地图', 'AdminKiAdminSetNextLayer');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (33, 2, 27, '更改游戏倍速', 'AdminSlomo');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (34, 2, 27, '玩家跳边', 'AdminForceTeamChange');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (35, 2, 27, '玩家降职', 'AdminDemoteCommander');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (36, 2, 27, '解散小队', 'AdminDisbandSquad');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (37, 2, 27, '把玩家踢出小队', 'AdminRemovePlayerFromSquad');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (38, 2, 27, '警告玩家', 'AdminWarn');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (39, 2, 27, '重开比赛', 'AdminRestartMatch');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (40, 2, 4, '刷新配置', 'refresh_config');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (41, 2, 4, '设置会员', 'vip');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (42, 2, 4, '封禁', 'ban');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (43, 2, 4, '分配权限', 'role');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (44, 1, 5, '玩家角色组', 'player_role');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (45, 1, 5, '面板角色组', 'amdin_role');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (46, 1, 5, '面板操作日志', 'admin_log');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (47, 1, 3, '删除', 'admin_delete');
INSERT INTO admin_auth (aa_id, aa_type, aa_pid, aa_name, aa_path) VALUES (48, 2, 2, '清除数据', 'data_clear');
INSERT INTO admin_role (ar_id, ar_name, ar_group) VALUES (1, '超级管理员', '*');`,
}
