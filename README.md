
### POST /signup
receives:

```json
{
    "email": email,
}
```
returns:
```json
{
    "email": signed email,
    "key": uuid key,
}
```

### POST /order
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
```json
{
    "acquired_wanted_item_amount": amount you have received,
    "surplus_given_item_amount": amount you returned to you, // if order is completed,
    "inorder_given_item_amount": amount still in order to execute, // else order isn't completed
}
```

## Inventory endpoints
prefixed by /$inventory_name
POST /additem
POST /delitem
GET /getprice

