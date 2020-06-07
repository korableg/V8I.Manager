[![Go Report Card](https://goreportcard.com/badge/github.com/korableg/OneCIBasesCreator)](https://goreportcard.com/report/github.com/korableg/OneCIBasesCreator)

<h3> Мини сервис по созданию файла со списком баз 1С (*v8i) на основании файла с настройками сервера (*lst)</h3>

<b>Параметры приложения:</b><br>

<b>-lst:</b> Путь к файлу 1CV8Clst.lst (по умолчанию расположен по адресу: %programfiles%\1cv8\srvinfo\reg_1541)<br>
<b>-ibases:</b> Путь к файлу iBases.v8i (по умолчанию расположен по адресу: %appdata%\1C\1CEStart)<br>

<b>Команды сервиса:</b><br>
<b>install:</b> Установить сервис<br>
<b>remove:</b> Удалить сервис<br>
<b>start:</b> Запустить сервис<br>
<b>stop:</b> Остановить сервис<br>
<b>pause:</b> Поставить сервис на паузу (Активный процесс скачивания будет работать пока не завершится)<br>
<b>continue:</b> Продолжить работу (после паузы)<br><br>
<b>Дополнительно для сервиса:</b><br>
<b>-instance:</b> Название сервиса (По умолчанию Downloader1C) (на случай если требуется развернуть несколько)<br>

lst и ibases могут быть перечислены через запятую, для обработки/создания сразу нескольких файлов<br>

Дерево групп информационных баз строится на основании поля "Описание", если данное поле не заполнено, то база укладывается в группу с названием сервера<br>
![Использование поля описание](https://github.com/korableg/OneCIBasesCreator/blob/master/blob/BaseProperties.png?raw=true)

Тестовый пример в репозитории генерирует такой список:<br>
![Список баз из тестового примера](https://github.com/korableg/OneCIBasesCreator/blob/master/blob/OneCStarter.png?raw=true)

Релизы подготавливаются традиционно под платформу Win64, если нужно скомпилировать под другие системы, напишите.

<b>Примеры:</b><br>
<b>Установка сервиса:</b> OneCIBasesCreator_Service install -lst "C:\Program Files\1cv8\srvinfo\reg_1541\1CV8Clst.lst" -ibases C:\Users\dmitry\AppData\Roaming\1C\1CEStart\ibases.v8i<br>
<b>Удаление сервиса:</b> OneCIBasesCreator_Service remove<br>
