# balance-api-exercise
the following configuration file is required     
### config.json     
```
{
	"PORT": <PORT NUMBER>,
	"DATA_SOURCE_NAME": <STRING WITH SQL DATA SOURCE NAME see https://en.wikipedia.org/wiki/Data_source_name>
}
```
### Overview
* /create-user [POST]    
accepts a JSON object with name to create a user.    
returns a JSON object with id of the created user   
```
$ curl --header "Content-Type: application/json" \
       --request POST \
       --data '{"name":"Andrew"}' \
       http://localhost:8000/create-user
     
$ {"result":{"uid":1}}
```
```
$ curl --header "Content-Type: application/json" \
       --request POST \
       --data '{"name":"David"}' \
       http://localhost:8000/create-user
     
$ {"result":{"uid":2}}
```

* /balance-operation [POST]    
either deposits or withdraws an amount of money from user's balance   
accepts JSON object with **uid, amount, note**.       
**amount** can be positive or negative, **note** is a description of the operation
```
$ curl --header "Content-Type: application/json" \
       --request POST \
       --data '{"uid":2, "amount": 1000.324, "note":"Random income"}' \
       http://localhost:8000/balance-operation
       
$ {"result":"ok"}
```

* /transaction [POST]    
transfers money from one user to another  
accepts JSON object with **from_uid, to_uid, amount, note**.       
**amount** must be positive, **note** is a description of the operation
```
$ curl --header "Content-Type: application/json" \
       --request POST \
       --data '{"from_uid":2, "to_uid":1,  "amount": 50.5, "note":"David to Andrew"}' \
       http://localhost:8000/transaction
       
$ {"result":"ok"}
```


* /notes [GET]    
returns a JSON object with notes for operations involving user's account       
accepts a query string with **Uid** parameter.
notes are sorted by timestamp
```
$ curl  http://localhost:8000/notes?Uid=2       
$ { 
    "result":[ 
               {"text":"Random income","timestamp":1601207819},
               {"text":"David to Andrew","timestamp":1601207989}
             ]
  }
```


* /balance [GET]    
returns a JSON object with balance value of user     
accepts a query string with **Uid** parameter.   
if an additional **Currency** parameter is present,        
return value will be converted to the specified currency   
```
$ curl "http://localhost:8000/balance?Uid=1&Currency="USD""      
$ {"result":0.64987224}
```




