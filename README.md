## A solution for [Wolt Summer 2024 Engineering Internships Task](https://github.com/woltapp/engineering-internship-2024) in Golang 
It is a simple web server written in **Golang** with only one POST REST API that calculates delivery fee based on user's input.

### How to run

```bash
git clone https://github.com/zianazari/delivery-fee-calculator.git
cd delivery-fee-calculator
go run .
```

### How to test with curl
```bash
curl -X POST "localhost:8080/calculate-delivery-fee" -H "Content-Type:application/json"  -d '{"cart_value":1000, "delivery_distance":1501, "number_of_items":3, "time":"2023-10-05T14:23:00Z"}'
```
#### Expected result 
```bash
{"delivery_fee":400}
```

### How to run all tests
```bash
go test .
```
#### Expected result 
```bash
ok      delivery_calculator_fee 0.485s
```
