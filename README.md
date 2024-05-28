# Сервис баннеров
В компании есть большое количество неоднородного контента, для которого необходимо иметь единую систему управления.  В частности, необходимо показывать разный контент пользователям в зависимости от их принадлежности к какой-либо группе. Данный контент мы будем предоставлять с помощью баннеров.
## Описание задачи
Необходимо реализовать сервис, который позволяет показывать пользователям баннеры, в зависимости от требуемой фичи и тега пользователя, а также управлять баннерами и связанными с ними тегами и фичами.
## Общие вводные
**Баннер** — это документ, описывающий какой-либо элемент пользовательского интерфейса. Технически баннер представляет собой JSON-документ неопределенной структуры. 
**Тег** — это сущность для обозначения группы пользователей; представляет собой число (ID тега). 
**Фича** — это домен или функциональность; представляет собой число (ID фичи).  
1. Один баннер может быть связан только с одной фичей и несколькими тегами
2. При этом один тег, как и одна фича, могут принадлежать разным баннерам одновременно
3. Фича и тег однозначно определяют баннер

Так как баннеры являются для пользователя вспомогательным функционалом, допускается, если пользователь в течение короткого срока будет получать устаревшую информацию.  При этом существует часть пользователей (порядка 10%), которым обязательно получать самую актуальную информацию. Для таких пользователей нужно предусмотреть механизм получения информации напрямую из БД.
## Условия
1. Используйте этот [API](https://github.com/avito-tech/backend-trainee-assignment-2024/blob/main/api.yaml)
2. Тегов и фичей небольшое количество (до 1000), RPS — 1k, SLI времени ответа — 50 мс, SLI успешности ответа — 99.99%
3. Для авторизации доступов должны использоваться 2 вида токенов: пользовательский и админский.  Получение баннера может происходить с помощью пользовательского или админского токена, а все остальные действия могут выполняться только с помощью админского токена.  
4. Реализуйте интеграционный или E2E-тест на сценарий получения баннера.
5. Если при получении баннера передан флаг use_last_revision, необходимо отдавать самую актуальную информацию.  В ином случае допускается передача информации, которая была актуальна 5 минут назад.
6. Баннеры могут быть временно выключены. Если баннер выключен, то обычные пользователи не должны его получать, при этом админы должны иметь к нему доступ.

## Дополнительные задания:
Эти задания не являются обязательными, но выполнение всех или части из них даст вам преимущество перед другими кандидатами. 
1. Адаптировать систему для значительного увеличения количества тегов и фичей, при котором допускается увеличение времени исполнения по редко запрашиваемым тегам и фичам
2. Провести нагрузочное тестирование полученного решения и приложить результаты тестирования к решению
3. Иногда получается так, что необходимо вернуться к одной из трех предыдущих версий баннера в связи с найденной ошибкой в логике, тексте и т.д.  Измените API таким образом, чтобы можно было просмотреть существующие версии баннера и выбрать подходящую версию
4. Добавить метод удаления баннеров по фиче или тегу, время ответа которого не должно превышать 100 мс, независимо от количества баннеров.  В связи с небольшим временем ответа метода, рекомендуется ознакомиться с механизмом выполнения отложенных действий 
5. Реализовать интеграционное или E2E-тестирование для остальных сценариев
6. Описать конфигурацию линтера

## Требования по стеку
- **Язык сервиса:** предпочтительным будет Go, при этом вы можете выбрать любой, удобный вам. 
- **База данных:** предпочтительной будет PostgreSQL, при этом вы можете выбрать любую, удобную вам. 
- Для **деплоя зависимостей и самого сервиса** рекомендуется использовать Docker и Docker Compose.


# Сервис баннеров

## Описание
Сервис баннеров предназначен для управления и отображения баннеров в зависимости от тегов пользователей и фичей. Баннеры представляют собой элементы пользовательского интерфейса в формате JSON.

## Установка и запуск
1. Убедитесь, что у вас установлены Docker и Docker Compose.
2. Склонируйте репозиторий: `git clone https://github.com/15Andrew43/backend-trainee-assignment-2024.git`
3. Перейдите в директорию проекта: `cd backend-trainee-assignment-2024`
4. Запустите сервис: `make stop && make && sleep 3 && bash ./db_init_queries/create_tables.sh`
5. Тут нужно создать какие-то теги чтобы все работало, 2 варианта
 - `bash ./db_init_queries/insert_test_data.sh`
 - `python3 ./db_init_queries/postgres/inserts.py n`
     - здесь n - сколько фич и тэгов мы хотим создать

Сервис будет доступен по адресу `http://localhost:8080`

## Использование
- Для получения баннера для пользователя выполните GET запрос к `/user_banner` с указанием тега и фичи пользователя.


`curl -X GET "http://localhost:8080/user_banner?tag_id=1&feature_id=1&use_last_revision=true" -H "token: AuthorizedUser"`

- Для управления баннерами (создание, обновление, удаление) выполните соответствующие запросы к `/banner` с использованием админского токена.


`curl -X POST "http://localhost:8080/banner" -H "token: Admin" -d '{"tag_ids":[1,2],"feature_id":1,"content":"{\"title\":\"New Banner\",\"text\":\"This is a new banner\",\"url\":\"https://example.com\"}","is_active":true}'`


`curl -X PATCH "http://localhost:8080/banner/1" -H "token: Admin" -d '{"tag_ids":[1,2, 3],"feature_id":2,"content":"{\"title\":\"Updated New Banner\",\"text\":\"Updated This is a new banner\",\"url\":\"https://Updated_example.com\"}","is_active":true}'`


`curl -X DELETE "http://localhost:8080/banner/1" -H "token: Admin"`

- Для получения всех баннеров с фильтрацией по тегам и фичам выполните GET запрос к `/banner`.

`curl -X GET "http://localhost:8080/banner?tag_id=1&feature_id=1&use_last_revision=true" -H "token: Admin"`

`curl -X GET "http://localhost:8080/banner?tag_id=1&use_last_revision=true" -H "token: Admin"`


## Тестирование
- Для запуска E2E тестов выполните: `make stop && make && sleep 3 && bash ./db_init_queries/create_tables.sh && bash ./db_init_queries/insert_test_data.sh && go test -count=1 ./test/E2E_test`.
- Для нагрузочного тестирования установите k6 и выполните соответствующие скрипты
     - `make stop && make && sleep 3 && bash ./db_init_queries/create_tables.sh && python3 ./db_init_queries/postgres/inserts.py 60000`
     - `k6 run ./test/load_test/create_banner.js`
     - `k6 run ./test/load_test/get_banner.js`
     - `k6 run ./test/load_test/update_banner.js`
     - `k6 run ./test/load_test/delete_banner.js`

Результаты нагрузочного теститрвоания можно посмотреть в папке `.test/load_test`

## Дополнительная информация
Для успешного выполенния bash-скриптов понадобятся `mongosh`, `psql`.

Для успешного нагрузочного тестирования необходим `k6`.

`docker logs my-go-server` - для просмотра логов