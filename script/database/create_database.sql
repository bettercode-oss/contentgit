CREATE DATABASE content_git WITH ENCODING = 'UTF8' LC_COLLATE = 'C.utf8' LC_CTYPE = 'C.utf8' TEMPLATE template0 IS_TEMPLATE = False;
CREATE USER content_git_app WITH ENCRYPTED PASSWORD '...';
GRANT ALL PRIVILEGES ON DATABASE content_git to content_git_app;
GRANT ALL ON SCHEMA public TO content_git_app;

-- https://github.com/tembo-io/pgmq
-- create the extension in the "pgmq" schema
CREATE EXTENSION pgmq;
-- creates the queue
SELECT pgmq.create('content');

GRANT ALL ON SCHEMA pgmq TO content_git_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON pgmq.q_content TO content_git_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON pgmq.a_content TO content_git_app;