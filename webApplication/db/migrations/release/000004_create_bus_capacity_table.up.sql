CREATE TABLE bus_capacity (
    id int NOT NULL AUTO_INCREMENT,
    bus_id int NOT NULL,
    route_id int NOT NULL,
    bus_cap int NOT NULL,
    bus_pass int NOT NULL,
    date_time datetime,
    PRIMARY KEY (`id`),
    UNIQUE KEY `BUS_ROUTE_INDX` (`bus_id`, `route_id`, `date_time`)
);