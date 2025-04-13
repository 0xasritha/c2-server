## Create Tasks 
```bash
curl -X POST localhost:8080/tasks \
-H "Content-Type: application/json" \
-d '[
{
"id": 0,
"type": "hello",
"options": {
"hello": "world",
"bye": "world"
}
},
{
"id": 1,
"type": "hehe",
"options": {
"haha": "hehe",
"ping": "pong"
}
}
]'
```