-- we don't know how to generate schema main (class Schema) :(
create table Name
(
	NameID INTEGER not null,
	Name TEXT not null,
	Sex int not null
)
;

create unique index Name_NameID_uindex
	on Name (NameID)
;

create unique index Name_Name_uindex
	on Name (Name)
;

-- unexpected locus for key
;

create table PredefinedAnswer
(
	PredefinedAnswerID INTEGER
		primary key
		 autoincrement,
	Category INTEGER not null,
	Answer TEXT not null
)
;

create table User
(
	UserID INTEGER
		primary key,
	Name TEXT not null,
	PasswordHash TEXT not null,
	SessionKey TEXT
)
;

create table Bot
(
	BotID INTEGER
		primary key
		 autoincrement,
	Name TEXT not null,
	Image TEXT not null,
	Gender INTEGER not null,
	User INTEGER not null
		constraint Bot_User_UserID_fk
			references User,
	Affection REAL not null,
	Mood REAL not null
)
;

create table Message
(
	MessageID INTEGER
		primary key
		 autoincrement,
	Bot INTEGER not null
		constraint Message_Bot_BotID_fk
			references Bot,
	Sender INTEGER not null,
	Timestamp DATETIME not null,
	Content TEXT not null,
	Rating REAL not null
)
;

create unique index User_Name_uindex
	on User (Name)
;

create unique index User_SessionKey_uindex
	on User (SessionKey)
;

create table Image
(
	ImageID INTEGER primary key autoincrement,
	Path TEXT not null,
)
;