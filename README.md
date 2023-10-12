# Minimalist Go Blog

Yet another personal blog writtein in GO. 
 
With HTMX templates and animations.

MySQL db and Go backend bundled with Docker. 



## Setup 

Edit /templates/*.html with your own resume / style


**Download and install**
- Go 
- MySQL
- Docker (optional)

## Download dependencies

gorilla/mux, mysql drivers, gosimple/slug, uuid, etc

`go mod download`


**env file example**

`PORT=80` if you plan to run this on your website

`MYSQL_HOST=any_hostname`

`MYSQL_PORT=3306`

`MYSQL_USER=any_username`

`MYSQL_PASSWORD=any_password`

`MYSQL_DATABASE=any_name`

`MYSQL_ROOT_PASSWORD=any_password2`
