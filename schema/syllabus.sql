CREATE TABLE lecture_rooms
(
    campus varchar(6),
    building varchar(20),
    room_id varchar(20),
    PRIMARY KEY(campus,building,room_id)
);

CREATE TABLE classroom_allocations
(
    lecture_id varchar(5),
    campus varchar(6),
    building varchar(20),
    room_id varchar(20),
    PRIMARY KEY(lecture_id,campus,building,room_id)
);

CREATE TABLE syllabus_base_infos
(
    year smallint,
    season varchar(8),
    day varchar(30),
    period varchar(30),
    teacher varchar(50),
    name varchar(100),
    lecture_id varchar(5),
    credits smallint,
    url varchar(100),
    type varchar(20),
    faculty varchar(4),
    PRIMARY KEY(lecture_id)
);