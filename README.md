##### disclaimer
This is some kind `hello world` project in `golang` for me.
So it may contain some not commonly used practices among `golang` developers (still learning, but for now: love error handling approach).

# Hancock
Service is meant to store and serve ad network SDKs performances. For long term storage is used MySql database which
need to be configured.

## Configuration

Service can be configured through configuration file or environment variables. Configuration file is located
at: `/cmd/config.json`. Variables from environment override configuration from `/cmd/config.json` file.

| Env Variable  | Json key  | Default value    |
|---|---|---|
| API_PORT | api->port  | 10000  |
| DB_HOST | db->host  | 127.0.0.1  |
| DB_PORT | db->port  | 3306  |
| DB_USERNAME | db->username  | root  |
| DB_PASSWORD | db->password  | root  |
| DB_DATABASE | db->database  | hancock  |

### User

For storing and requesting performances from the service you need to make http request thorough API. This request are
authenticated and authorized, so you need to authenticate yourself. By default, are available following testing users:

- user (username:`user`, password: `user`, role:`guest`)
- admin (username:`admin`, password: `admin`, role:`admin`)

## Api

Service has exposed api for user authentication and performances management.

### Auth

Authentication request:

- Endpoint: `/api/auth`
- Method: `POST`
- Body:

```json
{
  "username": "user",
  "password": "user"
}
```

Response:

- Status code: 200
- Cookies: `token`, `refresh_token`
- Body:

```json
{
  "token": "token_data",
  "refresh_token": "refresh_token_data"
}
```

Curl:

```bash
curl --location --request POST 'http://localhost:10000/api/auth' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username":"user",
    "password":"user"
}'
```

Re authentication request:

- Endpoint: `/api/auth`
- Method: `PUT`
- Cookies: `token`

Response:

- Status code: 200
- Cookies: `refresh_token`
- Body:

```json
{
  "token": "token_data"
}
```

Curl:

```bash
curl --location --request PUT 'http://localhost:10000/api/auth' \
--header 'Content-Type: application/json' \
--header 'Cookie: refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImNyZWF0ZWRfYXQiOiIyMDIxLTAxLTEwIDA3OjE1OjQ5IiwidXNlcm5hbWUiOiJ1c2VyIiwicm9sZSI6MTAsInNlc3Npb25faWQiOjIzLCJleHAiOjE2NDE4MTU0NTl9.ewkxcL36ketRAcvd11skGtGa-bdkN6H9akycu9q79zY; token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImNyZWF0ZWRfYXQiOiIyMDIxLTAxLTEwIDA3OjE1OjQ5IiwidXNlcm5hbWUiOiJ1c2VyIiwicm9sZSI6MTAsInNlc3Npb25faWQiOjIzLCJleHAiOjE2MTAyNzk4OTB9.5SnvELT-tV3xmrFBV8K8Z_XFYPuo8dUm5htj8zNqYqg'
```

De authentication request:

- Endpoint: `/api/auth`
- Method: `DELETE`
- Cookies: `token`

Response:

- Status code: 200

Curl:

```bash
curl --location --request DELETE 'http://localhost:10000/api/auth' \
--header 'Content-Type: application/json' \
--header 'Cookie: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImNyZWF0ZWRfYXQiOiIyMDIxLTAxLTEwIDA3OjE1OjQ5IiwidXNlcm5hbWUiOiJ1c2VyIiwicm9sZSI6MTAsInNlc3Npb25faWQiOjI2LCJleHAiOjE2MTAyODA0OTl9.pCp9cyRnm1tGnwjPHVYXuy_Mp022xuN8mKS6GGFN2GQ;' \
```

### Performances

Get performances:

- Endpoint: `/api/performances`
- Method: `GET`
- Cookies: `token`
- Query params: `country`,`platform`,`os_version`,`app_name`

Response:

- Status code: 200
- Body:

```json
{
  "country": "c1",
  "app": "a1",
  "banner": [
    {
      "sdk": "s1",
      "score": 2
    }
  ],
  "interstitial": [
    {
      "sdk": "s1",
      "score": 2
    }
  ],
  "reward": [
    {
      "sdk": "s1",
      "score": 2
    }
  ]
}
```

Curl

```bash
curl --location --request GET 'http://localhost:10000/api/performances?country=%22A%22&platform=%22B%22&os_version=2&app_name=B' \
--header 'Content-Type: application/json' \
--header 'Cookie: token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOjEsImNyZWF0ZWRfYXQiOiIyMDIxLTAxLTEwIDA3OjE1OjQ5IiwidXNlcm5hbWUiOiJ1c2VyIiwicm9sZSI6MTAsInNlc3Npb25faWQiOjI3LCJleHAiOjE2MTAyODA4OTN9.J22Kig1OR307Bn3x3qwUxLhDZ-m7mL5AF2WLSMb0NXo' 
```

Create performances (only users with admin user role allowed to create new performances):

- Endpoint: `/api/performances`
- Method: `POST`
- Cookies: `token`
- Body:

```json
[
  {
    "ad_type": "banner",
    "country": "c1",
    "app": "a1",
    "sdk": "s1",
    "score": 2
  },
  {
    "ad_type": "interstitial",
    "country": "c1",
    "app": "a1",
    "sdk": "s1",
    "score": 2
  },
  {
    "ad_type": "rewarded",
    "country": "c1",
    "app": "a1",
    "sdk": "s1",
    "score": 2
  }
]
```

Response:

- Status code: 201

## Make commands

### Local dockerized database

For development, we have `make` command with which you can create and destroy dockerized MySql database. Database has
configured following parameters:

- port: 3306
- username: root
- password: root
- database: hancock

Command for mannaging database are following:

```bash
make start-local-db
make stop-local-db
make restart-local-db
```

### App running

```bash
make run
```

### Unit testing

```bash
make test
```
### Base controlling commands

```bash
make login-user
make login-admin
```
 



