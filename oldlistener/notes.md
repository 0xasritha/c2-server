# Listener

### Results 
**Model**
```
Result {
   id uint32
   contents string 
   isSuccess bool  
}
```

**API Spec**
- `/results` (GET): 
- `/results` (POST): check if Results in request body are empty; if not, store each Result; get list of Tasks waiting to be served to implant, delete them from store so Tasks aren't served twice; send Tasks back to implant in response

---

### Tasks 
**Model** 
```
Task { 
    id   uint32
    type string
    options []string 
}
```

**API Spec** 
- `/tasks` (GET): get all Tasks  
- `/tasks` (POST): add UUID + store all Tasks in request body; copy tasks sent to implant into TaskHistory collection; return response that includes all Tasks that were added to the database 
---
### Task Record 

- Task Record API will return all tasks that were successfully served to implant + their associated results 

**Model** 
```
TaskRecord {
    // id is same as Task's ID? 
    task Task 
    results []Result
}
```
**API Spec**

- `/history` (GET): get all TaskRecords 

---
