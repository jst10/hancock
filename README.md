# Hancock
Project is meant to store and serve ad network SDKs performances.
For long term storage is used MySql database which need to be configured.

## Configuration
## Api
### Auth
Authentication request:
- Endpoint: `/api/auth`
- Method: `POST`
- Body:
```json
{
    "username":"user",
    "password":"user"
}
```

Response: 
- Status code: 200
- Cookies: `token`, `refresh_token`

Curl: 
```bash
curl --location --request POST 'http://localhost:10000/api/auth' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"user",
    "password":"user"
}'
```

### Performances

## Make commands
### Local dockerized database
### App running
### Base controlling commands
make login-user
make login-admin
## Sql database
