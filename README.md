### Keenetic API для версий 2.16+
На 3 не проверял пока, но думаю время придёт рано или поздно

### База:
```go
router := zyxel.NewZyxel("admin", "pass", "http://192.168.1.1")
err, cookies := router.Login() // Вообще он сам проверяет авторизацию при любом запросе и когда надо перезаходит
```

### Список девайсов:
```go
router := zyxel.NewZyxel("admin", "pass", "http://192.168.1.1")
err, cookies := router.Login() // Вообще он сам проверяет авторизацию при любом запросе и когда надо перезаходит
```

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