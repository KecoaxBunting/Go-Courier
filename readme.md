#Go-Courier

**Go-Courier** merupakan project microservices menggunakan gRPC

## Arsitektur Microservices

**1. Auth**
    **a. Auth_Server**: gRPC
    **b. Auth_Gateway**: gateway http
    **Route**
    1. POST `/register` registrasi user
    2. POST `/login` login user dan mendapatkan token
    3. GET `/api/me` mendapatkan user id

**2. Order**
    **a. Order_Server**: gRPC
    **b. Order_Gateway**: gateway http
    **Route**
    1. POST `/orders/` membuat order baru
    2. GET `/orders/` mendapatkan semua order
    3. GET `/orders/{orderId}` mendapatkan order dengan id tertentu
    4. PUT `/orders/{orderId}` update data order dengan id tertentu
    5. DELETE `/orders/{orderId}` delete data order dengan id tertentu

**3. Courier**
    **a. Courier_Server**: gRPC
    **b. Courier_Gateway**: gateway http
    **Route**
    1. POST `/couriers/` membuat courier baru
    2. GET `/couriers/` mendapatkan semua couriers
    3. GET `/couriers/{courierId}` mendapatkan courier dengan id tertentu
    4. PUT `/couriers/{courierId}` update data courier dengan id tertentu
    5. DELETE `/couriers/{courierId}` delete data courier dengan id tertentu
    6. PUT `/couriers/changeAvailability/{courierId}` ubah status available courier

**3. Delivery**
    **a. Delivery_Server**: gRPC
    **b. Delivery_Gateway**: gateway http
    **Route**
    1. POST `/assignCourier/` membuat delivery baru
    2. PUT `/completeOrder/` menyelesaikan delivery
    3. GET `/deliveries/` mendapatkan semua delivery
    4. GET `/deliveries/{deliveryId}` mendapatkan delivery dengan id tertentu
    5. DELETE `/deliveries/{deliveryId}` delete data delivery dengan id tertentu

## Note
- Sebelum menjalankan server dan gateway, tolong setup .env terlebih dahulu
- Jangan lupa build proto terlebih dahulu, gunakan command `make [nama_perintah]`
- Jalankan `command docker-compose up --build` untuk uji coba menggunakan docker

**David Adriel Alvyn**