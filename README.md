## Необходимый софт
1. Go
2. Docker
3. Docker Compose

## Программа может
1. Загружать мемы из альбома вездекода
2. Загружать мемы из других групп
3. Лайкать мемы
4. Продвигать мемы (они будут появляется чаще)
5. Составлять ленту из мемов

## Инструкция по использованию
0. Проверьте, что у вас не заняты порты: 8000, 5434
1. Запустите Postgres в docker контейнере
```bash
$ docker-compose up -d
```
2. Создайте нужные таблицы в БД
```bash
$ go run cmd/main.go db-up
```
3. Запустите программу в режиме сервера или в режиме сервера вместе с интерактивной консолью
```bash
$ go run cmd/main.go     # в режиме сервера
$ go run cmd/main.go cli # в режиме сервера и интерактивной CLI
```
Список доступных команд в режиме CLI
- help                          - выводит все доступные команды
- exit                          - выйти из терминала (учтите, что сервер перестанет работать)
- load_vezdekod                 - загрузить мемы из группы "Вездекод"
- load_group [ссылка на группу] - загрузить мемы из другой группы ( например https://vk.com/abstract_memes )
- feed                          - лента с мемами

## Endpoint'ы
- [GET] /api/memes/ - выводит все загруженые мемы<br><br>
- [POST] /api/memes/load-vezdekod - загрузит мемы из альбома Вездекода.<br><br> 
- [POST] /api/memes/load-group - загрузит мемы из вашей любимой группы ВК.  Запрос требует тело вида:
```json
{
  "password": "vezdekod",                      // Cтандартный пароль
  "address": "https://vk.com/informationmemes" // Ссылка на группу ВК
}
```

- [POST] /api/memes/feed - выводит мем <br>
Запрос требует тело вида:
```json
{
  "page": 1 // Номер мема по счету в ленте. Лента формируется основывая на дате поста и на продвигаемых постах
}
```

- [POST] /api/memes/like - лайкает мем от вашего IP адреса <br>
Запрос требует тело вида:
```json
{
  "meme_id": "15" // ID мема в БД
}
```
- [POST] /api/memes/promote - ставит мем на продвижение <br>
Запрос требует тело вида:
```json
{
  "meme_id": "15", // ID мемеа в БД
  "password": "vezdekod" // Стандартный пароль
}
```
