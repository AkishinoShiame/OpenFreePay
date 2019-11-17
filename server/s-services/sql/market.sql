DROP TABLE market;

CREATE TABLE market(
	order_id int not null primary key AUTO_INCREMENT, #會在我的出價用到
    date date not null default now(),
    user_id int not null,
    src_point_type text not null,
    dest_point_type text not null,
    revision_times int not null default 0,
    src_bid_points double not null,
    dest_ask_points double not null,
    username varchar(100) not null,
    status bool not null default false #訂單狀態
);