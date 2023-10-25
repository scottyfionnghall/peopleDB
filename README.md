# People Database




Web application based on REST API, made using Gin, GORM and PostgreSQL.


## About


This web app stores and returns database entries about different users, storing their name, age, gender and nationality.


## API Endpoints




- **`GET`** /username
- **`GET`** /username/all
- **`POST`** /username
- **`PATCH`** /username
- **`DELETE`** /username




### GET /username




HTTP Request to get information about a specific user


Available query params:




***At least one of them is required in the request***




- **id** - Find specific user based on their ID
- **name** - Find specific user based on their name, returns first similar entry
- **surname** - Find specific user based on their surname, returns first similar entry




**Request example:**


```shell
curl --location 'localhost:8080/user/?id=2'
```




### GET /username/all




HTTP request to get information about all users.




Available query params:


- **name** - Returns all users with entered name
- **surname** - Returns all users with entered surname
- **age** - Returns all users with entered age
- **gender** - Returns all users with entered gender
- **nationality** - Returns all users with entered nationality
- **order_by** - Defines in what order to return user entries, first word defines by what field to order entries, second one defines in what order (ascending or descending)


Example:
`id asc`,`name desc`,`surname asc` and etc.
Available fields by which to order response:
`id`,`name`,`age`,`surname`,`nationality`,`age`,`gender`.
- **page_size** - Defines page size; Default: 10
- **page** - Defines current page




**Request example:**


```shell
curl --location 'localhost:8080/user/all?name=Igor&page=1&page_size=10&orderby=name%20asc'
```




### POST /username




POST request to add users.


Body of the request must have a name and surname.


**Request example:**


```shell
curl --location 'localhost:8080/user' \
--header 'Content-Type: application/json' \
--data '{
"name":"Ben",
"surname":"Gilbert"
}'
```




### PATCH /username




PATCH request to change already existing entry. Accepts only one param which is user ID.


**Request example:**


```shell
curl --location --request PATCH 'localhost:8080/user?id=2' \
--header 'Content-Type: application/json' \
--data '{
"gender":"male"
}'
```




### DELETE /username




DELETE request to remove user by their ID.


**Request example:**


```shell
curl --location --request DELETE 'localhost:8080/user?id=1'
```
