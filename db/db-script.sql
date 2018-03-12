CREATE TABLE Bot
(
  BotID     INTEGER PRIMARY KEY AUTOINCREMENT ,
  Name      TEXT    NOT NULL,
  Image     TEXT    NOT NULL,
  Gender    TEXT    NOT NULL,
  User      INTEGER NOT NULL
    CONSTRAINT Bot_User_UserID_fk
    REFERENCES User,
  Affection REAL    NOT NULL,
  Mood      REAL    NOT NULL
);


CREATE TABLE Message
(
  MessageID INTEGER PRIMARY KEY AUTOINCREMENT,
  Bot       INTEGER  NOT NULL
    CONSTRAINT Message_Bot_BotID_fk
    REFERENCES Bot (BotID),
  Sender    INTEGER  NOT NULL,
  Timestamp DATETIME NOT NULL,
  Content   TEXT     NOT NULL,
  Rating    REAL     NOT NULL
);

CREATE TABLE PredefinedAnswer
(
  PredefinedAnswerID INTEGER PRIMARY KEY AUTOINCREMENT,
  Category           INTEGER NOT NULL,
  Answer             TEXT    NOT NULL
);

CREATE TABLE User
(
  UserID       INTEGER PRIMARY KEY AUTOINCREMENT,
  Name         TEXT NOT NULL,
  PasswordHash TEXT NOT NULL
);
