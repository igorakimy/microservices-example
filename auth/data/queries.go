package data

var (
	getAllUsersQuery = `
SELECT 
	id, 
	email, 
	first_name, 
	last_name, 
	password, 
	active,
	created_at,
	updated_at
FROM
    users
ORDER BY
    last_name
`

	getUserByEmailQuery = `
SELECT
	id,
	email,
    first_name,
    last_name,
    password,
    active,
    created_at,
    updated_at
FROM
    users
WHERE
    email = $1
`

	getUserByIDQuery = `
SELECT
	id,
	email,
    first_name,
    last_name,
    password,
    active,
    created_at,
    updated_at
FROM
    users
WHERE
    id = $1
`

	updateUserByIDQuery = `
UPDATE
	users
SET
    email = $1,
	first_name = $2,
	last_name = $3,
	active = $4,
	updated_at = $5
WHERE 
    id = $6
`

	deleteUserByIDQuery = `
DELETE
FROM
    users
WHERE 
    id = $1
`

	insertNewUserQuery = `
INSERT
INTO
	users(
		email,
	    first_name,
	    last_name,
	    password,
	    active,
	    created_at,
	    updated_at
	)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id
`

	updateUserPasswordQuery = `
UPDATE
	users
SET
    password = $1
WHERE
    id = $2
`
)
