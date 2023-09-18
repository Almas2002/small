## Getting Started
Вам нужен докер для старта этого проекта или локально установленные программы как postgres и jaeger.

### Installation

_

1. Clone the repo
   ```sh
   git https://github.com/Almas2002/small
   ```
2. Build golang project
   ```sh
   docker compose build
   ``` 
3. Вам будет нужны эти данные для подняние docker compose вы можете здесь изменить некоторые данные но не нельзя трогать sms данные
   ```dotenv
   PGPASSWORD=12345
   PGHOST=db
   IS_PRODUCTION=false
   POSTGRES_DB=small
   POSTGRES_USER=postgres
   POSTGRES_PORT=5432
   JAEGERHOSTPORT=jaeger:6831
   SMSHOST=465
   SMSEMAIL=almas.zhumakhanov@bk.ru
   SMSPASS=gBNrLu75CtCsh0QLVhEe
   SMSSMTP=smtp.mail.ru

   ```
4. Вам нужно поднять docker compose,здесь может быть некоторая ошибка в поднятии потому что postgres не успевает иницализироваться
 ```sh
   docker compose up -d
   ``` 

### Использование

1. Откройте ваш postman и перейдите по адресу http://localhost:8080/api/user/ здесь можно создать пользователя (POST) тело 
{
   "email":"hitba283@gmail.com"
   }

2. Откройте ваш postman и перейдите по адресу http://localhost:8080/api/product/ здесь можно создать продукт (POST) тело {
   "price":50.5,
   "title":"first"
   }

3. Откройте ваш postman и перейдите по адресу http://localhost:8080/api/product/:id здесь можно обновить товар по id (PUT) тело {
   "price":52.9
   }
4. Откройте ваш postman и перейдите по адресу http://localhost:8080/api/product/sub здесь можно подписаться на продукт (POST) тело {
   "user_id":1,
   "product_id":1
   }
5. Откройте ваш postman и перейдите по адресу http://localhost:8080/api/product/unsub здесь можно отписаться от продукта (DELETE) тело {
   "user_id":2,
   "product_id":1
   }
