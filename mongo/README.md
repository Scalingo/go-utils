# Package `mongo` v1.2.1

Useful tools around MongoDB.

```go
mongo.DefaultDatabaseName = "default-database"
mongo.Session(context.Background())
```

Create a new session based on `MONGO_URL`, connection will be initialized only
once, so you can call it every time your need a session.

Will wait until database is available.

If no `MONGO_URL` is defined, will use mongodb://localhost:27017/ + DefaultDatabaseName
