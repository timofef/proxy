# Proxy server
Домашнее задание по курсу "Анализ защищенности"

Макаров Тимофей АПО-31

Технопарк, 3 семестр
# Запуск
```
sudo docker build -t proxy . 
sudo docker run -p 8081:8081 -p 8082:8082 --name proxy -t proxy
```
# Использование
На порте 8081 - прокси сервер.

На порте 8082 - сервис для повтора запросов и проверки запросов на наличие XXE уязвимости.

/ - вывести список запросов.

/requests?id=<request_id> - вывести запрос с id, равным request_id.

/repeat?id=<request_id> - повтор запроса с id, равным request_id.

/scan?id=<request_id> - проверка запроса с id, равным request_id, на наличие XXE уязвимости.

Проверка на наличие XXE проводится с помощю подстаовки в запрос после строчки '''<?xml ...>'''
'''
<!DOCTYPE foo [
  <!ELEMENT foo ANY >
  <!ENTITY xxe SYSTEM "file:///etc/passwd" >]>
<foo>&xxe;</foo>
'''