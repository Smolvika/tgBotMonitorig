-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE allusers
(
    usersId    int(11) NOT NULL AUTO_INCREMENT,
    chatId     int(11) NOT NULL,
    useBotFunc tinyint(1) NOT NULL,
    PRIMARY KEY (usersId)
);
CREATE TABLE notif_change_currency (
                                       UserId int(11) NOT NULL AUTO_INCREMENT,
                                       ChatId int(11) NOT NULL,
                                       Сurrency VARCHAR(6) NOT NULL,
                                       Change_cost DECIMAL(7,2) DEFAULT NULL,
                                       PRIMARY KEY (UserId)
);
CREATE TABLE notif_currency (
                                UserId int(11) NOT NULL AUTO_INCREMENT,
                                ChatId int(11) NOT NULL,
                                Сurrency VARCHAR(6) NOT NULL,
                                PRIMARY KEY (UserId)
);
CREATE TABLE user_schat (
                            usersId int(11) NOT NULL AUTO_INCREMENT,
                            chatId VARCHAR(60) NOT NULL,
                            useBotFunc varchar(8) NOT NULL,
                            PRIMARY KEY (usersId)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE  allusers;
DROP TABLE  userschat;
DROP TABLE notif_change_currency;
DROP TABLE notif_currency;