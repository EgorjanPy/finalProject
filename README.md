# finalProject

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
Example curl request
```
curl http://127.0.0.1:8080/api/v1/calculatу
   -X POST -H "Content-Type: application/json"
   -d  "{\"expression\":\"2+2*2\"}"
```
Example response
```
result: 6.000000
```
If you want to run the calculator in the console, you can comment out the line app.RunServer() in main.go and uncomment // app.Run()
