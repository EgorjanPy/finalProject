# finalProject
Hello friend, this is a service for calculating arithmetic expressions
## Links
* [Bнструкция на русском языке](#установка)
* [Instructions in English](#installation)
## Установка
* Скачайте с официального сайта [golang](https://go.dev/dl/) на ваш компьютер если он не скачен
* Клонируйте репозиторий
```
git clone https://github.com/EgorjanPy/finalProject.git
```
* Перейдите в папку репозитория
```
cd ./finalProject/
```
## Запуск
Чтобы запустить, введите команду в терминал
```
go run ./cmd/ -port=8080
```
Вы можете выбрать другой порт запуска если хотите(по умолчанию 8080)
* Пример curl запроса
```
curl http://127.0.0.1:8080/api/v1/calculatу
   -X POST -H "Content-Type: application/json"
   -d  "{\"expression\":\"2+2*2\"}"
```
* Пример ответа
```
result: 6.000000
```
Если вы хотите запустить калькулятор в консоли, вы можете закомментировать строку app.RunServer() в main.go и раскомментировать // app.Run()

## Installation
 * Download from the official website [golang](https://go.dev/dl/) to your computer
 * Clone the repository
```
git clone https://github.com/EgorjanPy/finalProject.git
```
* Go to this folder
```
cd ./finalProject/
```
## Start
To run, enter the command
```
go run ./cmd/ -port=8080
```
You can choose another port if you want(default 8080).
* Example curl request:
```
curl http://127.0.0.1:8080/api/v1/calculatу
   -X POST -H "Content-Type: application/json"
   -d  "{\"expression\":\"2+2*2\"}"
```
* Example response
```
result: 6.000000
```
If you want to run the calculator in the console, you can comment out the line app.RunServer() in main.go and uncomment // app.Run()
