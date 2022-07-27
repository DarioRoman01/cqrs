# CQRS

## What is this? why?
this is a simple CRUD application on feeds the special thing in this project is the architecture wich is based on the CQRS architecture 
the main goal of this project was to learn more about micro services and how they work, the app have 3 services one in charge or read operations, other 
in charge of the write operations and other service that notify the client when a feed is created, deleted, or updated. for communicate this 
services i decided to use NATS as a message queue and nginx to map the requests to the needed service.

## Architecture overview
![architecture_diagram](https://cdn.discordapp.com/attachments/1001908460545388586/1001908493239992380/unknown.png)

## How to run it?
requirements:
```
docker 
docker-compose
```

in the root folder just run 
```bash
docker compose up -d --build
```

and you can try it out with tools like postman in the following url "localhost:8080"