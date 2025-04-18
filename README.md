# gifinator

fork of https://github.com/GoogleCloudPlatform/gifinator to run locally and use minio for storage instead of gcs.
look at TODO.md for what's left to be done.
look at ARCHITECTURE.md for how the code is structured.

# running (for development)
the following is what I do for development, you might hate it

## configure name resolution
add the following in `/etc/hosts`
```
127.0.0.1   minio
127.0.1.1   minio
```

## compile and docker compose
run `make all && docker compose build --no-cache && docker compose up --force-recreate --remove-orphans`

## visit web page
access http://localhost:8090 using a web browser, fill in the text, select one of the 3 options, click `Create` and wait

# running (for production)
you are not going to run this in production, let's be real

# documentation
env vars for each service are listed below

## redis
see [official docs](https://github.com/librenms/docker/blob/263c47e895850e6c7a4cafedd73fadd43b870711/doc/docker/environment-variables.md)

## minio
see [official docs](https://github.com/minio/minio/tree/9171d6ef651a852b48f39f828c3d01e30fbf4e9c/docs/config)

## render
| Variable       | Description                  | Default Value        |
|----------------|------------------------------|----------------------|
| `RENDER_PORT`  | port the service will run at | 8070                 |
| `MINIO_NAME`   | minio server domain name     | localhost            |
| `MINIO_PORT`   | minio server port number     | 9000                 |
| `MINIO_BUCKET` | minio bucket to be used      | gifbucket            |
| `MINIO_KEY`    | minio access key             | minioaccesskeyid     |
| `MINIO_SECRET` | minio secret key             | miniosecretaccesskey |

## gifcreator (server mode)
| Variable          | Description                         | Default Value        |
|-------------------|-------------------------------------|----------------------|
| `GIFCREATOR_PORT` | port the service will run at        | 8081                 |
| `SCENE_PATH`      | path to find files for gif creation | /scene               |
| `RENDER_NAME`     | renderer domain name                | localhost            |
| `RENDER_PORT`     | renderer port number                | 8070                 |
| `REDIS_NAME`      | redis server domain name            | localhost            |
| `REDIS_PORT`      | redis server port number            | 6379                 |
| `MINIO_NAME`      | minio server domain name            | localhost            |
| `MINIO_PORT`      | minio server port number            | 9000                 |
| `MINIO_BUCKET`    | minio bucket to be used             | gifbucket            |
| `MINIO_KEY`       | minio access key                    | minioaccesskeyid     |
| `MINIO_SECRET`    | minio secret key                    | miniosecretaccesskey |

## gifcreator (worker mode)
| Variable          | Description                         | Default Value        |
|-------------------|-------------------------------------|----------------------|
| `GIFCREATOR_PORT` | port the service will run at        | 8082                 |
| `SCENE_PATH`      | path to find files for gif creation | /scene               |
| `RENDER_NAME`     | renderer domain name                | localhost            |
| `RENDER_PORT`     | renderer port number                | 8070                 |
| `REDIS_NAME`      | redis server domain name            | localhost            |
| `REDIS_PORT`      | redis server port number            | 6379                 |
| `MINIO_NAME`      | minio server domain name            | localhost            |
| `MINIO_PORT`      | minio server port number            | 9000                 |
| `MINIO_BUCKET`    | minio bucket to be used             | gifbucket            |
| `MINIO_KEY`       | minio access key                    | minioaccesskeyid     |
| `MINIO_SECRET`    | minio secret key                    | miniosecretaccesskey |

## frontend
| Variable                 | Description                   | Default Value
|--------------------------|-------------------------------|--------------------
| `FRONTEND_PORT`          | port the service will run at  | 8090
| `FRONTEND_TEMPLATES_DIR` | directory for html templates  | /templates
| `GIFCREATOR_NAME`        | gifcreator server domain name | localhost
| `GIFCREATOR_PORT`        | gifcreator server port number | 8081

# legal stuff
the original is [here](https://github.com/GoogleCloudPlatform/gifinator) and its license and legal stuff apply, I'm not trying to steal anything
