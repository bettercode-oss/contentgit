CREATE DATABASE test_content_git WITH ENCODING = 'UTF8' LC_COLLATE = 'C.utf8' LC_CTYPE = 'C.utf8' TEMPLATE template0 IS_TEMPLATE = False;

\c test_content_git;
CREATE EXTENSION pgmq;
