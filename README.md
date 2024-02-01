# Learning RESTful API with Go

Simple RESTful API to create, read, update and delete books with MongoDB.

## Prepare Development Environment

```sh
docker compose -f docker-compose-dev.yml up -d
```

Now click on the "Attach to the running container" on VS Code.

```sh
go build
./go_restapi
```


## Endpoints

### GET

```txt
http://localhost:8000/books
```

### POST

URL

```txt
http://localhost:8000/books
```

Body

 ```txt
{
"isbn":"4545454",
"title":"First Three",
"author":{"firstname":"Harry",  "lastname":"White"}
}
 ```

Header:

```txt
Key:Value
Content-Type:application/json
```

### PUT

URL

```txt
http://localhost:8000/books/{specific-id}
```

Body

 ```txt
{
"isbn":"777777777777",
"title":"Updated Three",
"author":{"firstname":"Harry",  "lastname":"White"}
}
 ```

Header:

```txt
Key:Value
Content-Type:application/json
```

## DELETE

URL

```txt
http://localhost:8000/books/{id}
```

Acessing secure apis:

- Obtain a JWT token here: /static/authenticate.html
  - Enter `username:password > admin:12345,member:12345`
  - The response contains a JWT token for that program
- Use the token when calling any secure api:
  - set the Authorization request header and add the jwt token, like so:
  - Authorization: Bearer \<token\>
