## GOPAY Service for Go-Jek

The response look like this
```
{
    "Account": {
        "ID": 2,
        "CreatedAt": "2017-12-17T22:28:43.211754544+07:00",
        "UpdatedAt": "2017-12-17T22:28:43.211754544+07:00",
        "DeletedAt": null,
        "ExternalId": "1",
        "Type": "customer",
        "Amount": 0,
        "Passphrase": "password"
    },
    "Status": "OK"
}
```

POST '/' response Status = "OK"/"FAILED"
GET '/' response Status = "OK"/"UNAUTHORIZED"
PUT '/' response Status = "OK"/"INSUFFICIENT"/"UNAUTHORIZED"
PUT 'topup' Status = "OK"/"INVALID"/"UNAUTHORIZED"