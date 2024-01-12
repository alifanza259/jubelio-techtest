<h1>Simple Chat App</h1>

# How to run
1. Run `docker run --name jubelio-db -d -p 5123:5432 -e POSTGRES_PASSWORD=password postgres`
2. Populate the database with data using sql script: 
- Run `docker cp ./migrate.sql jubelio-db:/migrate.sql`
- Run `docker exec -u postgres jubelio-db psql postgres postgres -f /migrate.sql`
3. Run `go mod tidy`
4. Run `go run main.go`
5. Start testing

Example testing: 
- Open localhost:8081/public/homepage.html with Google Chrome
- Login with email: a@gmail.com, password: pass
- Click go to chat app
- Click 'Open Chat', insert id 2
- Type a message, click send
- Open another tab, open localhost:8081/public/homepage.html, login with email: b@gmail.com, password: pass
- Click go to chat app
- Click 'Open Chat', insert id 1
- See previous message, type another message and click send

# Notes
Currently there are 3 users for login after running migration: 
1. Email: a@gmail.com, Password: pass, Name: Andi
2. Email: b@gmail.com, Password: pass, Name: Bakri 
3. Email: c@gmail.com, Password: pass, Name: Ciko

# Endpoint List:
Frontend:
- GET /public/homepage.html -> Show login form & redirect to chat page
- GET /public/chat.html -> Show chat page

Backend:
- GET /login -> Login with email and password (Check notes, currently there is no register feature)
- GET /messages -> Get messages between 2 users
- POST /search -> Search messages in local database
- POST /messages -> Send messages from front end to supabase
- GET /ws -> Accept websocket connection
