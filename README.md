## GOPAY Service for Go-Jek

The response look like this
```
{
    "account": {
        "ID": 2,
        "CreatedAt": "2017-12-17T22:28:43.211754544+07:00",
        "UpdatedAt": "2017-12-17T22:28:43.211754544+07:00",
        "DeletedAt": null,
        "external_id": "1",
        "type": "customer",
        "amount": 0,
        "passphrase": "<password>"
    },
    "status": "OK"
}
```

POST '/' response status = "OK"/"FAILED"
GET '/' response status = "OK"/"UNAUTHORIZED"
PUT '/' response status = "OK"/"INSUFFICIENT"/"UNAUTHORIZED"
PUT 'topup' status = "OK"/"INVALID"/"UNAUTHORIZED"