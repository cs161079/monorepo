ALTER TABLE `schedulemaster` ADD COLUMN `sdc_days` VARCHAR(7);
ALTER TABLE `schedulemaster` ADD COLUMN `sdc_months` VARCHAR(12);
alter table route modify route_descr varchar(150);
alter table route modify route_descr_eng varchar(150);