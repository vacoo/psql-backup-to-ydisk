# POSTGRESQL BACKUP TO YANDEX DISK

Утилита для создания, сохранения и восстановления резервных копий из базы данных postgresql. Хранилище Yandex Disk

В установленные сроки по cron заданию выполняется бэкап базы данных. Затем этот бэкап загружается в Yandex Disk. В диске резервные копии находятся в папке Приложения/имя-вашего-проекта (`https://disk.yandex.ru`). Копии разделены по месяцам чтобы удобно было ориентироватся. Если кое где будет ошибка будет отправлен отчет в указанную почту.

## Получение токена доступа
1. Региструем свое приложение в yandex `https://oauth.yandex.ru/client/new`. Выбираем Яндекс.Диск REST API и права: Доступ к информации о диске, Доступ к папке приложения на диске. Сохраняем полученные данные.
2. Заходим по этой ссылке `https://oauth.yandex.ru/authorize?response_type=token&client_id=<client_id>` и получаем токен доступа. Токен действует 1 год.
3. Полученный YANDEX_DISK_ACCESS_TOKEN & YANDEX_DISK_APP_FOLDER записываем в конфигурацию.

## Пример конфигурации в docker-compose.yml
docker-compose.example.yml

## Бэкап по требованию

`docker exec -it <container_name> sh /home/backup`

## Восстановление из бэкапа

1. `docker exec -it <container_name> sh /home/app restore 2020-01/backup_2020-01-17_15-10.gz`
2. Бэкап будет сохранен в /home/backups

## Использование

MAIL_SMPT_HOST=smtp.yandex.ru MAIL_SMPT_PORT=465 MAIL_SMPT_USER=it-alarm@yandex.ru MAIL_SMPT_PASS=jitgfftikelmripv MAIL_TO=it-alarm@yandex.ru PROJECT=GogoTaxi YANDEX_DISK_ACCESS_TOKEN=AgAAAAAMLrD-AAYUPZxeHoLACkiYnQnsk-DnUwQ go run main.go
