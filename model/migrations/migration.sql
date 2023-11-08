CREATE TABLE users (
    id serial PRIMARY KEY,
    tenant_id int DEFAULT 0,
    email varchar(100) NOT NULL UNIQUE,
    password varchar(100) NOT NULL,
    is_active boolean DEFAULT true,
    plan_id int DEFAULT NULL,
    is_google_signin boolean DEFAULT false NOT NULL,
    create_date timestamp DEFAULT current_timestamp,
    update_date timestamp DEFAULT current_timestamp,
    is_verified boolean DEFAULT false NOT NULL,
    show_product_tour boolean DEFAULT true NOT NULL,
    trial_start_date timestamp DEFAULT NULL,
    trial_end_date timestamp DEFAULT NULL
)

CREATE TABLE sessions (
    id serial PRIMARY KEY,
    session_id char(36) NOT NULL UNIQUE,
    email varchar(100) NOT NULL,
    location varchar(100) NOT NULL DEFAULT 'uae',
    exp bigint NOT NULL,
    session_start_time timestamp DEFAULT current_timestamp,
    session_end_time timestamp DEFAULT NULL,
    CONSTRAINT fk_email
        FOREIGN KEY(email)
            REFERENCES users(email)
)

CREATE TABLE api_key_info (
    id serial PRIMARY KEY,
    secret_key char(32) NOT NULL UNIQUE,
    owner_email varchar(100) NOT NULL,
    domain varchar(100) NOT NULL,
    CONSTRAINT fk_email
        FOREIGN KEY(owner_email)
            REFERENCES users(email)
)

CREATE TABLE api_key_permissions (
    id serial PRIMARY KEY,
    secret_key char(32) NOT NULL UNIQUE,
    permitted_end_point varchar(100),
    CONSTRAINT fk_secret_key
        FOREIGN KEY(secret_key)
            REFERENCES api_key_info(secret_key)
)