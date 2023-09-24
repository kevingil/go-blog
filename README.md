# Go Blog

Personal blog and CMS


**Notes**

Default templates use Tailwind CSS

`restart: unless-stopped` in Docker compose


## Setup Instructions


Edit /views/*.gohtml with your own resume / style

**Download and install**
- Go 
- MySQL
- Docker


**env file example**

`PORT=8080`

`MYSQL_HOST=any_hostname`

`MYSQL_PORT=3306`

`MYSQL_USER=any_username`

`MYSQL_PASSWORD=any_password`

`MYSQL_DATABASE=any_name`

`MYSQL_ROOT_PASSWORD=any_password2`


**Docker Build**


From a Droplet

Clone from Github



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





