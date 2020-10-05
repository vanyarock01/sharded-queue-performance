# sharded-queue-performance

Нагрузочное тестирование шардированной очереди на базе Tarantool + Cartridge с использованием Pandora + Yandex.Tank

## Инструменты

Понадобятся:

- [Docker](https://docs.docker.com/engine/install/) 18.XX+
- [Tarantool](https://www.tarantool.io/en/download/)
- [cartridge-cli](https://github.com/tarantool/cartridge-cli#installation)
- опционально [Golang](https://golang.org/doc/install) 1.13+

## Собираем пушку

1. Если go есть в системе:
    ```shell
    go get github.com/tarantool/go-tarantool \
           github.com/spf13/afero            \
           github.com/yandex/pandora
    
    GOOS=linux GOARCH=amd64 go build tnt_queue_gun.go
    ```
2. Если go нет в системе, но есть докер:
    ```shell
    ./build.sh
    ```

## Собираем приложение

```shell
$ cd queue-app
$ cartridge build
```

## Запускаем приложения и настраиваем кластер

Для запуска инстансов кластера:

```shell
$ cartridge start -d
```

Конфигурацию кластера выполним скриптом bootstrap.lua:
```shell
$ tarantool bootstrap.lua
```

В случае успеха, по адресу [localhost:8081](localhost:8081) в браузере будет видная следующая конфигурация кластера:
![](./media/cluster.png)

## Создаем очереди для тестирования

В терминах sharded-queue экземпляр очереди - труба (tube). Создать ее можно, как через бинарное апи, так и через конфигурацию кластера. Воспользуемся вторым методом.

Перейдем во вкладку [code](http://localhost:8081/admin/cluster/code) в левом меню веб морды.
В редакторе кода создадим файл `tubes.yml` в который поместим конфигурацию интересующей нас очереди:
```yaml
test-tube:
    temporary: false
    driver: sharded_queue.drivers.fifo
```
![](./media/create_tube.png)

Остается нажать кнопку *Apply* и дождаться сообщения об успешном выполнении операции.

## Запускаем нагрузочные тесты

```shell
$ cd ..

$ docker run -v $(pwd):/var/loadtest      \
           -v $SSH_AUTH_SOCK:/ssh-agent \
           -e SSH_AUTH_SOCK=/ssh-agent  \
           --net host                   \
           -it direvius/yandex-tank
```

**NOTE** Для docker for mac может понадобиться заменить `localhost` на алиас `host.docker.internal` в файле [tnt_queue_load.yaml](./tnt_queue_load.yaml).


Более детальную информацию про установку и настройку Yandex.Tank можно найти в [документации](https://yandextank.readthedocs.io/en/latest/install.html).

