CREATE TABLE users (
    id uuid primary key not null,
    email character varying(100),
    password text,
    admin boolean,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE settings (
    id uuid primary key not null,
    refresh_rate integer,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE feeds (
    id uuid primary key not null,
    title text,
    url text unique not null,
    description text,
    site_url text,
    favicon text,
    domain text,
    added_by_id uuid references users(id),
    last_fetch timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE INDEX ON feeds ((lower(url)));

CREATE TABLE items (
    id uuid primary key not null,
    feed_id uuid references feeds(id),
    title text,
    author text,
    description text,
    link text,
    posted_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

CREATE TABLE subscriptions (
    id uuid primary key not null,
    user_id uuid not null,
    feed_id uuid not null,
    unique (user_id, feed_id)
);

CREATE TABLE item_meta (
    item_id uuid references items(id),
    user_id uuid references users(id),
    saved boolean,
    read_time timestamp with time zone,
    unique (item_id, user_id)
);

CREATE TABLE subscription_categories (
    subscription_id uuid not null,
    category_id uuid not null,
    feed_id uuid not null,
    unique (subscription_id, category_id, feed_id)
);
