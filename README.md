# Personal Blog


>Minimalist resume Go blog with mysql db and htmx frontend


## Setup 

Edit /templates/*.html with your own resume / style


**Download and install**
- Go 
- Docker (optional)

## Download dependencies

gorilla/mux, mysql drivers, gosimple/slug, uuid, etc

`go mod download`


.env example

```sh
PORT=80 #if you plan to run this on your website
PROD_DSN=your_database_key #or you can set it up with Docker
```



## Screenshots


### Dashboard

![main](https://cdn.kevingil.com/dashboard-main.png)


### Editor

![editor](https://cdn.kevingil.com/dashboard-editor.png)



