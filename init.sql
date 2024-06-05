SET ROLE TO videos;

SELECT current_user;

CREATE TABLE IF NOT EXISTS videos (
    video_id TEXT,
    trending_date TEXT,
    title TEXT,
    channel_title TEXT,
    category_id INTEGER,
    publish_time TEXT,
    tags TEXT,
    views INTEGER,
    likes INTEGER,
    dislikes INTEGER,
    comment_count INTEGER,
    thumbnail_link TEXT,
    comments_disabled BOOLEAN,
    ratings_disabled BOOLEAN,
    video_error_or_removed BOOLEAN,
    description TEXT
);

\copy videos FROM '/URLShortener/data.csv' WITH (FORMAT csv, HEADER);