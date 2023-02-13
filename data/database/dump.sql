\c postgres

CREATE TABLE movies (
    id      int, 
    title   varchar(250), 
    rating  int, 
    notes   text,
    year    varchar(50),
    PRIMARY KEY(id)
);

CREATE TABLE metadata (
    id int,
    movie_id int,
    poster varchar(250),
    release_date varchar(50),
    PRIMARY KEY(id),
    CONSTRAINT fk_movie
        FOREIGN KEY(movie_id)
        REFERENCES movies(id)
        ON DELETE CASCADE
);