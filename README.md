# helloworld-go

#Endpoints

- get a message for a username
```bash
curl http://localhost:8081/hello/{username}
```

- update or create a new user
```bash
curl -X PUT -I -d '{"dateOfBirth":"2003-10-12"}'  http://localhost:8081/hello/{username}
```


To run it locally postgres should be configured and schema should be created manually

```sql
CREATE TABLE public.users
(
    username text COLLATE pg_catalog."default" NOT NULL,
    dateofbirth date NOT NULL,
    CONSTRAINT "User_pkey" PRIMARY KEY (username)
)
WITH (
    OIDS = FALSE
)
```