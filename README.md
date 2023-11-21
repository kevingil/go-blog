# Personal Blog


>Minimalist resume Go blog with mysql db and htmx frontend


## Setup 

Edit /templates/*.html with your own resume / style


**Download and install**
- Go 
- MySQL
- Docker (optional)

## Download dependencies

gorilla/mux, mysql drivers, gosimple/slug, uuid, etc

`go mod download`


.env example

```sh
PORT=80 #if you plan to run this on your website
MYSQL_HOST=any_hostname
MYSQL_PORT=3306
MYSQL_USER=any_username
MYSQL_PASSWORD=any_password
MYSQL_DATABASE=any_name
MYSQL_ROOT_PASSWORD=any_password2
```



## Screenshots


Dashboard

![dashboards-articles.png](https://cdn.kevingil.com/dashboard-articles.png)

![dashboards-profile.png](https://cdn.kevingil.com/dashboard-profile.png)


