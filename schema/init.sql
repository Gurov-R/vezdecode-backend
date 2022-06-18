CREATE TABLE memes
(
    id serial not null unique,
    vk_id bigint not null unique,
    image_url text not null,
    timestamp_ bigint not null,
    likes_count int not null,
    promoted boolean not null default false
);

CREATE TABLE likes
(
    id serial not null unique,
    meme_id int not null,
    user_ip varchar(255) not null
);