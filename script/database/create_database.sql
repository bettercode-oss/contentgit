-- create the database
CREATE DATABASE content_git WITH ENCODING = 'UTF8' LC_COLLATE = 'C.utf8' LC_CTYPE = 'C.utf8' TEMPLATE template0 IS_TEMPLATE = False;

-- switch to the database
\c content_git;

-- create the extension in the "pgmq" schema
CREATE EXTENSION pgmq;

-- creates the queue
SELECT pgmq.create('content');