# Puskesmas

1. Add `.env` containing

```sh
SERVER_PORT=8002
DB_HOST=dbhost
DB_USERNAME=dbusername
DB_PASSWORD=dbpassword
DB_NAME=dbhost
ADMIN_USERNAME=
ADMIN_PASSWORD=
SERVICE_ACCOUNT_FILE=
```

2. Install all dependencies on `landing`

```sh
npm i
```

3. Add firebase service account to `service-account.json` file

4. Build

```sh
./run prod
```
