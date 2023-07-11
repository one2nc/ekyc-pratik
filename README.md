# EKyc

The purpose of this exercise is for you to get familiar with building HTTP API with tests on a real world use case. We won’t build a ton of requirements or features, but this exercise will give you enough idea about building REST APIs in Golang that use Database, Async workers and Caches, which are the most common components of any web application.
## Features

- Upload Image
- Face match score generation
- OCR for extracting data from ID cards

## How to use
- Clone this repository [https://github.com/one2nc/ekyc-pratik.git]
- Run [make setup] to setup docker containers of postgres and minio from docker compose file
- Run [make run] to build and start server


## Migrations
- Migration scripts are in db/migrations folder
- Migrations are run using golang-migrate package
- Migrations run during initialization of server
```sh
# command to add a migration file
migrate create -ext sql -dir <directory_path> -seq <migration_name>
```

## Enviroment Variables
- Refer env.example file to create your own env file

| ENV_VARIABLE | Description |
| :-------- | :------------------------- |
| `DB_HOST` | Host for databse connection 
| `DB_USER` | Databse user
| `DB_PASSWORD` | Database password
| `DB_PORT` | Database port
| `DB_NAME` | Database name
| `SERVER_PORT` | Server port
| `SERVER_HOST` | Server Host
| `DB_MIGRATION_FILE` | DB migration file path
| `MINIO_ACCESS_KEY` | Access key for minio server
| `MINIO_SECRET_KEY` | Secret key for minio server
| `MINIO_IMAGE_BUCKET_NAME` | Name of image bucket
| `MINIO_IMAGE_ENDPOINT` | Minio endpoint for api

## Minio Setup
- MinIO is an open-source, self-hosted object storage server that is compatible with Amazon S3 cloud storage service.
- Images uploaded by customer are stored into minio server
- Minio runs in an container using docker compose file
- Console can be accessed by going to {localhost:9001}
- Login to console using MINIO_ROOT_USER and MINIO_ROOT_PASSWORD provided in docker compose file
- After logging in generate access and secret key, this will be used in env variables [MINIO_ACCESS_KEY,MINIO_SECRET_KEY]
- Create a bucket for images and set bucket name into [MINIO_IMAGE_BUCKET_NAME] env variable


## API Reference
- Postman collection is in  [ekyc.postman_collection.json] file

#### Signup

```http
  POST /api/v1/auth/singup
```
| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `name` | `string` | **Required**. |
| `email` | `string` | **Required**. |
| `plan` | `string` | **Required**. |

#### Image Upload

```http
  POST /api/v1/image/upload
```
| Headers | Description     
| :-------- | :------- |
| `Access-Key` |  **Required**. |
| `Secret-Key` |  **Required**. |

| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `image` | `file` | **Required**. |
| `image_type` | `string` | **Required**. |

