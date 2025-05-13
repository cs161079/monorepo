CREATE TABLE IF NOT EXISTS opswCronRuns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    runtime     datetime NOT NULL,
    finishtime datetime
);