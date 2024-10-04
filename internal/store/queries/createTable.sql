create table if not exists prices (
    id integer primary key unique,
    name text not null,
    price integer not null,
    from_date date not null
)