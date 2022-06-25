# Lofty users API

> Hosted online on
> https://calm-meadow-78917.herokuapp.com/

To run locally, ensure all ENV variables are set in `.env` file, see `sample.env`

```bash
make server
```

## Endpoints

- **URL**

  `/users` - returns All users

- **Method:**

  `GET`

- **URL Params**

  **Optional:**

  `email=[alphanumeric]`
  `page=[numeric]`
  `per_page=[numeric]`
  `date_of_birth_to[alphanumeric]`
  `date_of_birth_from[alphanumeric]`

- **Success Response:**

* **Code:** 200 <br />

```
https://calm-meadow-78917.herokuapp.com/users?per_page=2

{
  "data": [
    {
      "id": 1,
      "email": "rhollier0@un.org",
      "gender": "2022-04-22T16:20:40Z",
      "last_name": "Hollier",
      "first_name": "Roxanne",
      "date_of_birth": "Female"
    },
    {
      "id": 2,
      "email": "lguerola1@opera.com",
      "gender": "2022-05-10T19:45:14Z",
      "last_name": "Guerola",
      "first_name": "Lonee",
      "date_of_birth": "Female"
    }
  ],
  "count": 2,
  "current_page": 1
}
```

- **URL**

  `/users/{email}` - returns specific user by email

- **Method:**

  `GET`

- **Success Response:**

* **Code:** 200 <br />

```
https://calm-meadow-78917.herokuapp.com/users/lguerola1@opera.com

{
  "data": {
    "id": 2,
    "email": "lguerola1@opera.com",
    "gender": "2022-05-10T19:45:14Z",
    "last_name": "Guerola",
    "first_name": "Lonee",
    "date_of_birth": "Female"
  }
}
```

- **Error Response:**

* **Code:** 404 <br />

```
https://calm-meadow-78917.herokuapp.com/users/rando

{
  "status": 404,
  "message": "User not found"
}
```
