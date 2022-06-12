Use json. Endpoints returns BadRequest 400 for invalid requests. bodies.

### POST /signup
---
receives:

```json
{
    "email": email,
}
```
returns:
OK 200
```json
{
    "email": signed email,
    "key": uuid key,
}
```
or 
Conflict 409, email already exists.

### POST /order
---
Orders are partially fulfillable by design. Individual txs are indicated by the terminal or sent to your email.

receives:
```json
{
    "email": the email you signed up with,
    "key": your account key,
    "given_item_inventory": inventory name,
    "given_item": name of given item,
    "given_item_amount": amount,
    "wanted_item_inventory": name of inventory,
    "wanted_item": name of wanted item,
    "wanted_item_amount": amount,
}
```
returns: 
Created 201
```json
{
    "acquired_wanted_item_amount": amount you have received,
    // if order is completed,=
    "surplus_given_item_amount": amount you returned to you, 
    // else order isn't completed
    "inorder_given_item_amount": amount still in order to execute,
}
```
or
NotAcceptable 406, maximum given item cost < minimum wanted item cost
or
BadRequest 400, invalid (email, key) pair.

# Inventory endpoints
prefixed by /$inventory_name

### POST /add
---
Add an item.
receives:
```json
{
    "name": item name,
    "price_min": float,
    "price_max": float,
}
```
returns:
Created 201
or 
Conflict 409
or
BadRequest 400, didn't satisfy logical checks


### DELETE /del
---
Delete existing item.
receives:
```json
{
    "name": item name,
}
```
returns:
NoContent 204, succesfully deleted
or
BadRequest 400, no such item
or
ExpectationFailed 417, wasn't deleted, server error.

### PUT /update
---
Update price of existing item.
receives:
```json
{
    "name": item name,
    "price_min": float,
    "price_max": float,
}
```
returns:
OK 200
or
BadRequest 400, no such item
or
ExpectationFailed 417, wasn't updated, server error.


### GET /list
---
List all items in the inventory.
receives empty request body.
returns:
OK 200
```json
{
    {item},
    .
    .
    .
}
```
or
InternalServerError 500, server error.

### GET /check
---
Check inventory for the item
receives:
```json
{
    "item_name": string,
}
```
returns:
Conflict 409, exists
or
NotFound 404, not exists.

### GET /cost
---
Calculate cost of wanted items. (wanted_item_price * wanted_amount)
receives:
```json
{
    "wanted_item": string,
    "wanted_amount": integer,
}
```
returns:
OK 200
```json
{
    "cost": float,
}
```
or
BadRequest 400, no such item.


### GET /calculate
---
Calculate how many wanted items can be acquired with the budget. (budget / wanted_item_price)
receives:
```json
{
    "budget": float,
    "wanted_item": string,
}
```

returns:
OK 200
```json
{
    "result": float,
}
```
or
BadRequest 400, wanted item not found.

