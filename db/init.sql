USE ${MYSQL_DATABASE};

source /docker-entrypoint-initdb.d/blog.sql;
