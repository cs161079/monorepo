ALTER TABLE `schedulemaster` ADD COLUMN `sdc_days` VARCHAR(7);
ALTER TABLE `schedulemaster` ADD COLUMN `sdc_months` VARCHAR(12);
ALTER TABLE `route` MODIFY `route_descr` VARCHAR(150);
ALTER TABLE `route` MODIFY `route_descr_eng` VARCHAR(150);