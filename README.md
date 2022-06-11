Orders are partially fulfillable by design.

### POST /signup
returns:
```json
{
    key: uuid_key,
}
```

### POST /order
returns: 
```json
{
    acquired_wanted_item_amount: x,
    surplus_given_item_amount: y, // if order is completed,
    inorder_given_item_amount: z, // elif order isn't completed
}
```

POST /{inventory_name}/additem
POST /{inventory_name}/delitem
GET /{inventory_name}/getprice

