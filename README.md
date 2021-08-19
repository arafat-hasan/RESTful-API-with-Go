 Learning RESTful API with Go


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
Key:Value
Content-Type:application/json


## PUT

URL
```
http://localhost:8000/books/{random-id}
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
Key:Value
Content-Type:application/json

## DELETE

URL
```
http://localhost:8000/books/{id}
```

