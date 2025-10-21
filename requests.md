# Примеры запросов:

GET /health - тестовый запрос с отображением статуса 
    
    GET http://localhost:8080/health

    
POST /tasks (создание) - создать новый task с заданным названием. id назначается сам.
   
    POST http://localhost:8080/tasks 
        Header Content-Type: application/json"
        Body {"title":"Купить молоко"}

        
GET /tasks (список) - вывести список все ранее созданных task'ов
 
    GET http://localhost:8080/tasks
    GET http://localhost:8080/tasks?q=молоко

    
GET /tasks/{id} - вывести конкретный task по id
  
    GET http://localhost:8080/tasks/1
