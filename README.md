[![Go Report Card](https://goreportcard.com/badge/github.com/korableg/OneCIBasesCreator)](https://goreportcard.com/report/github.com/korableg/OneCIBasesCreator)

# Cервис по созданию файла со списком баз 1С (*v8i) на основании файла с настройками сервера (*lst)

## Параметры приложения:

### Команды сервиса:
- **install**: Установить как службу Windows
- **remove**: Удалить службу
- **version**: Вывести в консоль версию приложения

### Флаги
- **--lst**: Путь к файлу 1CV8Clst.lst (по умолчанию расположен по адресу: %programfiles%\1cv8\srvinfo\reg_1541). Можно ввести несколько значений через запятую.
- **--v8i**: Путь к файлу iBases.v8i (по умолчанию расположен по адресу: %appdata%\1C\1CEStart). Можно ввести несколько значений через запятую.
- **--cfg**: Путь к конфигу приложения в формате .yaml [Пример](https://github.com/korableg/V8I.Manager/blob/master/assets/config_example.yaml). Если одновременно указан и конфиг файл и lst и v8i, то в приоритете будут они.
- **--help**: Вывести в консоль справку

### Описание
Дерево групп информационных баз строится на основании поля "Описание", если данное поле не заполнено, то база укладывается в группу с названием сервера  
![Использование поля описание](https://github.com/korableg/OneCIBasesCreator/blob/master/assets/BaseProperties.png?raw=true)

Тестовый пример в репозитории генерирует такой список:  
![Список баз из тестового примера](https://github.com/korableg/OneCIBasesCreator/blob/master/assets/OneCStarter.png?raw=true)

### Пример установки сервиса
`v8imanager.exe install --cfg C:\Users\User\config.yaml`
