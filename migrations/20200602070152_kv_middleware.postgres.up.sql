CREATE TABLE keys 
(
    id SERIAL,
    key VARCHAR(150),
    value TEXT,
    type VARCHAR(10),
    create_time TIMESTAMP default current_timestamp,
    update_time TIMESTAMP,
    created_by INT,
    approved_by INT default 0,
    status INT,
    PRIMARY KEY (id)
);

CREATE TABLE users
(
    id SERIAL,
    username VARCHAR(150),
    token VARCHAR(150),
    email VARCHAR(150),
    create_time TIMESTAMP default current_timestamp,
    created_by INT,
    status INT,
    PRIMARY KEY (id)
);

CREATE TABLE roles 
(
    id SERIAL,
    prefix VARCHAR(150),
    permission VARCHAR(150),
    create_time TIMESTAMP default current_timestamp,
    created_by INT,
    status INT,
    PRIMARY KEY (id)
);

CREATE TABLE user_access 
(
    user_id INT,
    role_id INT,
    create_time TIMESTAMP default current_timestamp,
    created_by INT,
    update_time TIMESTAMP,
    updated_by INT,
    status INT,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE canary_keys
(
    key_id INT,
    ip VARCHAR(20),
    status INT,
    PRIMARY KEY (key_id, ip)
);

INSERT INTO users(username, email, token, status, created_by) VALUES('admin', 'admin@tokopedia.com', '', 1, 0);
INSERT INTO users(username, email, token, status, created_by) VALUES('bb', 'bb@tokopedia.com', '', 1, 0);
INSERT INTO users(username, email, token, status, created_by) VALUES('user', 'user@tokopedia.com', '', 1, 0);

INSERT INTO roles(prefix, permission, status, created_by) VALUES('service', 'admin', 1, 0);
INSERT INTO roles(prefix, permission, status, created_by) VALUES('service/risk/sauron', 'lead', 1, 0);
INSERT INTO roles(prefix, permission, status, created_by) VALUES('service/risk/sauron', 'user', 1, 0);
INSERT INTO roles(prefix, permission, status, created_by) VALUES('service/risk/sauron', 'superuser', 1, 0);

INSERT INTO user_access(user_id, role_id, status, created_by) VALUES(1, 1, 1, 0);
INSERT INTO user_access(user_id, role_id, status, created_by) VALUES(2, 2, 1, 0);
INSERT INTO user_access(user_id, role_id, status, created_by) VALUES(3, 3, 1, 0);