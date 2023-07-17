# EKyc
EKyc is a system that manages user kyc. It provides API based solution for face matching and OCR.<br/> 
The purpose of this exercise is for you to get familiar with building HTTP API with tests on a real-world use case. This exercise will give you enough idea about building REST APIs in Golang that uses Database, Async workers and Caches, which are the most common components of any web application.
----
## Features

- Customer Signup. 
- Upload Image.
- Face match score generation.
- OCR for extracting data from ID cards.
----
## How to use
- Clone this repository.
```sh
git clone https://github.com/one2nc/ekyc-pratik.git
```
----

## How to setup
- To set up docker containers of postgres, redis and MinIO from docker-compose file, run.

```sh
make setup
```
- To take down docker containers, run.
```sh
make setup-down
```
----

## MinIO Setup
- MinIO is an open-source, self-hosted object storage server compatible with Amazon S3 cloud storage service.
- Images uploaded by customers are stored on MinIO server.
- MinIO runs in a container using docker-compose file.
- MinIO admin console can be accessed by opening [localhost:9001](http://localhost:9001) in the browser.
- Login to the console using `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD` provided in docker-compose file.
- After logging in generate `access and secret key`, which will be used in env variables `MINIO_ACCESS_KEY, MINIO_SECRET_KEY`.
  - On the left, there is a menu, click on **Access Key**.
  - Now click on **Create Access Key**  and then click on **Create**.
  - Copy newly generated access and secret key and set them into the above-mentioned env variables.
  - <img width="400" alt="Screenshot 2023-07-13 at 1 01 36 PM" src="https://github.com/one2nc/ekyc-pratik/assets/46452464/01217b50-61dd-4400-8b4a-a318e4e6ce3e">

- Create a bucket for images and set the bucket name into `MINIO_IMAGE_BUCKET_NAME` env variable.

  - On the left, there is a menu, click on **Buckets**.
  - Now click on **Create Bucket**, Entar bucket name in input and click on **Create**.
  - Set this bucket name in the `MINIO_IMAGE_BUCKET_NAME` env variable.
  - <img width="400" alt="Screenshot 2023-07-13 at 1 02 35 PM" src="https://github.com/one2nc/ekyc-pratik/assets/46452464/21dfb238-ed05-4693-8b9a-cbbe097df8bb">
- MinIO API's can be accessed through `localhost:9000`. This has to be set in the `MINIO_IMAGE_ENDPOINT` env variable.
----
## Enviroment Variables
- Refer `.env.example` file to create your own `.env` file in the root of the project.

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
| `MINIO_ACCESS_KEY` | Access key for MinIO server
| `MINIO_SECRET_KEY` | Secret key for MinIO server
| `MINIO_IMAGE_BUCKET_NAME` | Name of image bucket
| `MINIO_IMAGE_ENDPOINT` | MinIO endpoint for api
| `REDIS_ADDRESS` | Redis address
| `REDIS_PORT` | Redis port

----
## Migrations
- Migration scripts are in the `db/migrations` folder.
- Migrations are run using [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4) package.
- Migrations run during the initialisation of server.
```sh
# command to add a migration file
migrate create -ext sql -dir <directory_path> -seq <migration_name>
```
----

## How to Run

- To build and start the server, run.
```sh
make run
```
----
## API Reference
#### Postman Setup
- [Postman](https://www.postman.com/) is an API platform for building and using APIs.
- Postman collection is in  `ekyc.postman_collection.json` file in the root of the project.
- Open postman and import this collection. After importing you can see all the requests under `ekyc` collection.
- Environment variables:

| Postman ENV_VARIABLE | Description |
| :-------- | :------------------------- |
| `baseUrl` | set by default to http://127.0.0.1:3000 using pre request script .
| `access_key` | set using script by extratcting values from response of the signup api.
| `secret_key` | set using script by extratcting values from response of the signup api.

**Note**: First you need to create and select an environment.

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
| `access_key` | `string` |  
| `secret_key` | `string` |  

---

#### Image Upload

```http
  POST /api/v1/image/upload
```
| Headers | Description     
| :-------- | :------- |
| `Access-Key` |  **Required**. |
| `Secret-Key` |  **Required**. |
- These headers are set automatically by env variables and scripts.

| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `image` | `file` | **Required**. |
| `image_type` | `string` | **Required**|
- Valid image_type values are `(face or id_card)`.

#### response
| Body Parameters | Type     | 
| :-------- | :------- | 
| `image_id` | `string` |  

---
#### Face Match

```http
  POST /api/v1/image/face-match
```
| Headers | Description     
| :-------- | :------- |
| `Access-Key` |  **Required**. |
| `Secret-Key` |  **Required**. |
- These headers are set automatically by env variables and scripts.


| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `image_id_1` | `string,uuid` | **Required**. |
| `image_id_2` | `string,uuid` | **Required**. |
- Upload images using upload image api. After uploading you will get image_id for each image in response.
- Use these image id's as image_id_1 and image_id_2 in the face match api. Note that both ids should not be same.

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
- These headers are set automatically by env variables and scripts.


| Body Parameters | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `image_id` | `string,uuid` | **Required**. |
- Upload images using upload image api. After uploading you will get image_id for each image in response.
- Use this image id in the ocr api body to get data. Note that provided image should be of type id_card.

#### response
| Body Parameters | Type     |
| :-------- | :------- | 
| `data.name` | `string` |  
| `data.dob` | `string` |
| `data.gender` | `string` |
| `data.address` | `string` |
| `data.pincode` | `string` |
| `data.idNumber` | `string` |
