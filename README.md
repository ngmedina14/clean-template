# clean-architecture

```
.
├── cmd //  the entry point of your application.
├── internal //  contains the core logic of your application, including the API handlers, controllers, services, and repositories.
│   ├── api // contains the functions that handle incoming HTTP requests and send responses. (standalone)
│   ├── controller //  often contains structs that bundle together the dependencies required by your handlers. (Method on struct)
│   ├── database //  contains the code that sets up and manages the connection to your database.
│   ├── repository // contains the code that interacts directly with the database.(low-level module)
│   └── service // contains the business logic of your application.(high-level module)
└── schema // contains sql data for migration
```