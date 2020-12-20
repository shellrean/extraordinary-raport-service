# Extraordinary-raport-service

Extraordinary-raport-service adalah aplikasi yang melayani management penilaian siswa (sikap dan perkembangan capaian) dengan adanya aplikasi ini rekapitulasi dan statistika menjadi lebih mudah. Sebagai service tentu saja ini tidak hanya dipakai untuk 1 aplikasi client tetapi dapat menjadi data utama dari banyak aplikasi nantinya

## Installation
Clone aplikasi kemudian jalankan
``` bash
$ git clone https://github.com/shellrean/extraordinary-raport-service.git
$ cd extraordinary-raport-service 
$ go run app/main.go
```

## Configuration
Create ``config.yml`` file
``` yaml
# Server config
server:
  host: 127.0.0.1
  port: 9000
  timeout:
    server: 30
    read: 15
    write: 10
    idle: 5
# Database config
database:
  username: "postgres"
  password: ""
  dbname: ""
  host: "localhost"
  port: "5432"
  timezone: "Asia/Jakarta"
# Context config
context:
  timeout: 2
# Release config
release: true
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)