# Puskesmas

1. Add `.env` containing
```sh
SERVER_PORT=8002
DB_HOST=dbhost
DB_USERNAME=dbusername
DB_PASSWORD=dbpassword
DB_NAME=dbhost
```

2. Install all dependencies on `admin`, `dist`, and `landing`

```sh
npm i
```

3. Build
```sh
./build.sh
```