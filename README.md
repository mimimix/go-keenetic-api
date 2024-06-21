### Keenetic API для версии 2.16
Ниже не проверялись. 3.0+ тоже не проверял, но думаю время придёт рано или поздно

```bash
go get -u github.com/tucnak/telebot
```

### Аунтефикация:
```go
router := zyxel.NewZyxel("admin", "pass", "http://192.168.1.1")
err, cookies := router.Login() // Вообще он сам проверяет аунтефикация при любом запросе и когда надо перезаходит
```

### Список девайсов:
```go
router := zyxel.NewZyxel("admin", "pass", "http://192.168.1.1")
devices, err := router.DeviceList() // Получает список устройств, поля смотреть в автокомплите
```

---

### Пуллер клиетов
Сообщает когда в сеть заходит или выходит клиент

```go
router := zyxel.NewZyxel("admin", "pass", "http://192.168.1.1")
poll := poller.NewPoller(router, 5) // 5 - это интервал
poll.Start()
poll.Stop()

go func() {
	for {
		event := <-poll.Channel
		s, _ := prettyjson.Marshal(event)
		fmt.Println(string(s))
	}
}()
```

В канал кидает событие с девайсом и новым состоянием онлайна

```go
type PollEvent struct {
	IsOnline bool
	Device   *zyxel.Device
}
```