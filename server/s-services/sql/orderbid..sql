DROP TABLE orderbid;

CREATE TABLE orderbid(
	c_id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
    date DATE NOT NULL DEFAULT NOW(),
	order_id int NOT NULL,
    user_id int NOT NULL,
    username varchar(100) not null, #show that name
    target_src_bid_points double NOT NULL,
    target_dest_ask_points double NOT NULL,
    transcation_state bool NOT NULL DEFAULT false
)