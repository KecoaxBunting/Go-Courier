# Go-Courier

**Go-Courier** merupakan project microservices menggunakan gRPC

## Arsitektur Microservices

## 1. Auth
a. **Auth_Server**: gRPC <br>
b. **Auth_Gateway**: gateway http

**Route**
1. POST `/register` registrasi user <br>
2. POST `/login` login user dan mendapatkan token <br>
3. GET `/api/me` mendapatkan user id <br>

## 2. Order
a. **Order_Server**: gRPC <br>
b. **Order_Gateway**: gateway http
    
**Route**
1. POST `/orders/` membuat order baru <br>
2. GET `/orders/` mendapatkan semua order <br>
3. GET `/orders/{orderId}` mendapatkan order dengan id tertentu <br>
4. PUT `/orders/{orderId}` update data order dengan id tertentu <br>
5. DELETE `/orders/{orderId}` delete data order dengan id tertentu <br>

## 3. Courier
a. **Courier_Server**: gRPC <br>
b. **Courier_Gateway**: gateway http

**Route**
1. POST `/couriers/` membuat courier baru <br>
2. GET `/couriers/` mendapatkan semua couriers <br>
3. GET `/couriers/{courierId}` mendapatkan courier dengan id tertentu <br>
4. PUT `/couriers/{courierId}` update data courier dengan id tertentu <br>
5. DELETE `/couriers/{courierId}` delete data courier dengan id tertentu <br>
6. PUT `/couriers/changeAvailability/{courierId}` ubah status available courier <br>

## 4. Delivery
a. **Delivery_Server**: gRPC <br>
b. **Delivery_Gateway**: gateway http

**Route**
1. POST `/assignCourier/` membuat delivery baru <br>
2. PUT `/completeOrder/` menyelesaikan delivery <br>
3. GET `/deliveries/` mendapatkan semua delivery <br>
4. GET `/deliveries/{deliveryId}` mendapatkan delivery dengan id tertentu <br>
5. DELETE `/deliveries/{deliveryId}` delete data delivery dengan id tertentu <br>

## Note
- Sebelum menjalankan server dan gateway, tolong setup .env terlebih dahulu
- Jangan lupa build proto terlebih dahulu, buat folder proto, gunakan command `make [nama_perintah]` (bisa lihat di MakeFile)
- Jalankan `command docker-compose up --build` untuk uji coba menggunakan docker

**David Adriel Alvyn**