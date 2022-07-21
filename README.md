# directusapi

todo:

API methods

- [x] Authenticate
- [x] List Item
- [x] Insert
- [x] Create (partials)
- [x] Get By ID
- [x] Delete
- [x] Set Item
- [x] Update (partials)

Types support

- [x] generic id
- [x] enumeration string constants
- [x] string as primary key
- [x] float
- [x] time
- [x] boolean
- [x] pointers
- [x] array
- [x] object (there are known issues, I was only able to make it working with map[string]string)
- [] reference
- [] array of objects (repeater)

Error handling

- [] missing required input field for create/insert/update/set

Testing

- [] e2e tests working locally
- [] e2e tests working in CI
- [] compatibility with directus v9
- [] insert vs set

Batch operations

- [] batch update
- [] batch insert
- [] batch delete

Other

- [] godoc
- [] fileupload
- [] embeded structs as W or R
