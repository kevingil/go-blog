# Go Blog

Personal blog and CMS


## Dependencies
`"github.com/go-sql-driver/mysql"`

`"github.com/joho/godotenv"`

`"github.com/gorilla/mux"`


## Setup Instructions


Edit /views/*.gohtml with your own resume / style

**Download and install**
- Go 
- MySQL
- Docker


**env file example**

`PORT=80` if you plan to run this on your website

`MYSQL_HOST=any_hostname`

`MYSQL_PORT=3306`

`MYSQL_USER=any_username`

`MYSQL_PASSWORD=any_password`

`MYSQL_DATABASE=any_name`

`MYSQL_ROOT_PASSWORD=any_password2`


**Docker Build**


*I hosted mine in a DigitalOcean droplet for $$$ savings.*


Clone repo, install dependencies, then:


Initialize database

`docker-compose build`

`docker exec -it blog-db /bin/sh`

`mysql -u root -p < init.sql`


Then you can restart the build

`docker-compose down`

`docker-compose build`

`docker-compose up -d`

Blog should be serving on `localhost:8080`

Register a user at /register


**Troubleshooting**

If you get this error

`dial tcp XXX.XX.X.X:XXXX: connect: connection refused`

Just restart the app

`docker-compose restart app`

It's a problem with my Dockerfile, in my TODO. 




