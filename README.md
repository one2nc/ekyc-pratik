# EKyc
Ekyc is a system that manages user kyc. It provides API based solution for face matching and OCR 
The purpose of this exercise is for you to get familiar with building HTTP API with tests on a real-world use case. This exercise will give you enough idea about building REST APIs in Golang that use Databases, Async workers and Caches, which are the most common components of any web application.

## Features

- Signup
- Upload Image
- Face match score generation
- OCR for extracting data from ID cards

## How to use
- Clone this repository 
```sh
git clone https://github.com/one2nc/ekyc-pratik.git
```
- To setup docker containers of postgres and minio from docker-compose file, run

```sh
make setup
```
- To build and start the server, run
```sh
make run
```


## Migrations
- Migration scripts are in the `db/migrations` folder
- Migrations are run using [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4) package
- Migrations run during initialisation of server
```sh
# command to add a migration file
migrate create -ext sql -dir <directory_path> -seq <migration_name>
```

## Enviroment Variables
- Refer `.env.example file` to create your own `.env` file in the root of the project

| ENV_VARIABLE | Description |
| :-------- | :------------------------- |
| `DB_HOST` | Host for database connection 
| `DB_USER` | Database user
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
- MinIO is an open-source, self-hosted object storage server compatible with Amazon S3 cloud storage service.
- Images uploaded by customers are stored on minio server
- Minio runs in a container using docker-compose file
- Minio admin console can be accessed by opening [localhost:9001](http://localhost:9001) in the browser
- Login to the console using `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD` provided in docker-compose file
- After logging in generate access and secret key, this will be used in env variables `MINIO_ACCESS_KEY,MINIO_SECRET_KEY`
- Create a bucket for images and set the bucket name into `MINIO_IMAGE_BUCKET_NAME` env variable


## API Reference
#### Postman Setup
- [Postman](https://www.postman.com/) is an API platform for building and using APIs
- Postman collection is in  `ekyc.postman_collection.json` file in the root of the project
- Open postman and import this collection. after importing you can see all the requests under `ekyc` collection
- Setup environment and add these variables `baseURL` `access_key` `secret_key`. these variables will be used while making the request

#### Signup

```http
  POST /api/v1/auth/singup
```
| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `name` | `string` | **Required**. |
| `email` | `string` | **Required**. |
| `plan` | `string` | **Required**. |
#### response
| Body Parameters | Type     | 
| :-------- | :------- |
| `accessKey` | `string` |  
| `secretKey` | `string` |  

---

#### Image Upload

```http
  POST /api/v1/image/upload
```
| Headers | Description     
| :-------- | :------- |
| `Access-Key` |  **Required**. |
| `Secret-Key` |  **Required**. |
- Get Access-Key and Secret-Key values from signup api

| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `image` | `file` | **Required**. |
| `image_type` | `string` | **Required**|
- Valid image_type values are (face or id_card)

#### response
| Body Parameters | Type     | 
| :-------- | :------- | 
| `imageId` | `string` |  

---
#### Face Match

```http
  POST /api/v1/image/face-match
```
| Headers | Description     
| :-------- | :------- |
| `Access-Key` |  **Required**. |
| `Secret-Key` |  **Required**. |
- Get Access-Key and Secret-Key values from signup api

| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `image_id_1` | `string,uuid` | **Required**. |
| `image_id_2` | `string,uuid` | **Required**. |
- Get image_id_1 and image_id_2 value from Image Upload api

#### response
| Body Parameters | Type|
| :-------- | :------- |
| `score` | `int` |  |
---
#### OCR 
```http
  POST /api/v1/image/ocr
```
| Headers | Description     
| :-------- | :------- |
| `Access-Key` |  **Required**. |
| `Secret-Key` |  **Required**. |

| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `image_id` | `string,uuid` | **Required**. |
- Get image_id Image Upload api

#### response
| Body Parameters | Type     |
| :-------- | :------- | 
| `data.name` | `string` |  
| `data.dob` | `string` |
| `data.gender` | `string` |
| `data.address` | `string` |
| `data.pincode` | `string` |
| `data.idNumber` | `string` |
