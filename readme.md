# Welcome to my Music API REST!

Hi! In this doc you will learn to use the api rest step by step.

This api have 7 routes, each with his explication more down.

## Routes

For users have 5 routes.

- /api/v1/users **GET**: function for get all users in all the system
- /api/v1/users/{id} **GET**: function for get a user by his id
- /api/v1/users/{id} **DELETE**: function for delete a user by his id

For the auth exist 2 routes.

- /api/v1/register **POST**: endpoint for register a user. It require 3 parameters in the request body:
  `{ "name": "David Vargas Hernandez", "email": "david@gmail.com", "password": "secretPassword" }`

- /api/v1/login **POST**: endpoint for login a user in the system. It require 3 parameters in the request body:
  `{ "name": "David Vargas Hernandez", "email": "david@gmail.com", "password": "secretPassword" }`
  After it return a token like: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6IiIsImV4cCI6MTY2NTcxMjAxNH0.1q4e2iHSJlgbu3fkAPdRSAHyp-y2FscKjIhLH-0c8P0`
  NOTE: If you is not register you will can't login, first you register.

For music and songs have 2 routes.

- /api/v1/songs **GET**: endpoint for get all the songs in the data base.
- /api/v1/music **POST**: endpoint for get the song for search from iTunes API for beet the params. Require a _term_ for the search, for example:
  `http://localhost:3000/api/v1/music?term=luis+miguel+ayer`
  It will search and automatic save the result in the database
  NOTE: If you is not register o is don't login, you will can’t watch or save a song, first you login.

## Technologies

This api his made in **Golang** with **Mux** para el CRUD, **Postgres** as database and others librarys for encrypt password, generate the token for the auth and **Gorm** for generate the tables in the database from the models created.
