 Learning RESTful API with Go


Simple RESTful API to create, read, update and delete books with MongoDB.


```
go build
./go_restapi

```


# Endpoints
## GET

```
http://localhost:8000/books
```

 ## POST

URL
```
http://localhost:8000/books
```
Body
 ```
{
"isbn":"4545454",
"title":"First Three",
"author":{"firstname":"Harry",  "lastname":"White"}
}
 ```

Header:
```
Key:Value
Content-Type:application/json
```

## PUT

URL
```
http://localhost:8000/books/{specific-id}
```
Body
 ```
{
"isbn":"777777777777",
"title":"Updated Three",
"author":{"firstname":"Harry",  "lastname":"White"}
}
 ```

Header:
```
Key:Value
Content-Type:application/json
```

## DELETE

URL
```
http://localhost:8000/books/{id}
```


Acessing secure apis:

  - Obtain a JWT token here: /static/authenticate.html
    - Enter `username:password > admin:12345,member:12345`
    - The response contains a JWT token for that program

 - Use the token when calling any secure api:
    - set the Authorization request header and add the jwt token, like so:
    - Authorization: Bearer \<token\>

