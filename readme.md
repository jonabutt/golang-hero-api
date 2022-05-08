## GOLANG HERO API

This is a side project to learn golang. This API was created using golang without no external libraries. Features include list all heros, get hero, create hero, delete hero and update hero the data is stored in memory.

Created by Jonathan Buttigieg, email: jonathanbuttigieg1@gmail.com.

### Prerequisites

To run this project you need:

[GoLang](https://go.dev/doc/install/)

### Installing and running

Clone the project:

```
git clone https://github.com/jonabutt/golang-hero-api
```

Run the application:

```
go run ./main/main.go
```

### Testing the API

These are example of http requests using CURL.

List all heros:

```
curl -X GET localhost:8081/heros
```

Get hero:

```
curl -X GET localhost:8081/heros/1
```

Create a hero:

```
curl -X POST --data '{"id":"2","name":"Catwoman","firstName":"Selina","secondName":"Kyle","place":"Gotham"}' localhost:8081/heros
```

Delete a hero:

```
curl -X DELETE localhost:8081/heros/1
```

Update a hero:

```
curl -X PUT --data '{"id":"1","name":"Catwoman","firstName":"Selina","secondName":"Kyle","place":"Gotham"}' localhost:8081/heros
```
