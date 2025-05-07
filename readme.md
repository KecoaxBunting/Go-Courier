#Go-Courier

**Go-Courier** merupakan project microservices menggunakan gRPC

## Arsitektur Microservices

- **Auth_Server**: gRPC login dan registrasi user
- **Auth_Gateway**: gateway untuk login dan registrasi user

- **Order_Server**: gRPC insert dan view order
- **Order_Gateway**: gateway untuk insert dan view order

- **Courier_Server**: gRPC insert dan view kurir
- **Courier_Gateway**: gateway untuk insert dan view kurir

## Note
- Sebelum menjalankan server dan gateway, tolong setup .env terlebih dahulu
- Jangan lupa build proto terlebih dahulu, gunakan command `make [nama_perintah]`
- Jalankan `command docker-compose up --build` untuk uji coba menggunakan docker

**David Adriel Alvyn**