package main

/*
	Реализовывать протокол HTTP с нуля не будем, маршрутизацию запросов тоже,
	возьмём готовые.
	Хранилище - полностью самодельное, даже реализация связанного списка.

	Реализация - в чистом виде "proof-of-concept", даже контроль входных данных
	отсутствует.

	Есть опции командной строки, для настройки адреса и порта для простушивания,
	а так же для изменения TTL записей.

	Общая идея такова.

	Клиенту доступно три метода: PUT (или POST, здесь эффект будет одинаковый),
	GET и DELETE.

	Для упрощения реализации (чтобы не городить "JSON RPC", поскольку вопрос
	не в этом), ключ и значение задаются прямо в URI плейнтекстом. В случае,
	когда потребуется включать пробелы или спецсимволы, моджно кодировать
	в base64, например.

	Таким образом, запрос вида
	curl -X PUT http://localhost:8881/storage/mykey/somevalue
	добавляет в хранилище ключ "mykey" со значением "somevalue", возвращает
	код 200 и текст "OK".

	Запрос вида
	curl -X GET http://localhost:8881/storage/mykey
	вернёт либо код 200 и "somevalue", либо код 400 и "FAILURE", если такого
	ключа ещё/уже нет.

	Запрос вида
	curl -X DELETE http://localhost:8881/storage/mykey
	удалит ключ "mykey" и вернёт код 200 с текстом "OK", либо вернёт код 400
	с текстом "FAILURE", если такого ключа в хранилище нет.


	Механика хранилища.

	Хранилище организовано в виде связанного списка, работающего в качестве очереди.

	Для ускорения выборки существует ассоциативный массив (map), со ссылками по
	ключам на узлы данных с соответсвующими значениями (оверхед, но выигрыш в
	скорости тем больше, чем больше записей в хранилище).

	Узлы бывают двух типов: узел данных и "филлер".

	Механизм устаревания последовательно извлекает узлы из головы очереди.
	Если извлечён узел данных - то он отбрасывается в связи с устареванием.
	(Удаляются все ссылки на него, после чего объект снесёт при "сборке
	мусора".) Если извлечён "филлер", то механизм устаревания "засыпает" на
	одну секунду.

	Каждую секунду в хвост очереди автоматически добавляется новый филлер.

	При создании новой записи создаётся узел данных и добавляется в хвост очереди.
	При обновлении существующей записи, её узел данных уничтожается и
	создаётся заново в хвосте очереди (эффективней было бы просто править
	ссылки, но это усложняет и запутывает код, при том, что не требуется в
	данной задаче).

	Таким образом, из хранилища каждую секунду удаляются записи, устаревшие	на
	данный момент, без перебора всего хранилища каждый раз.

	Для синхронизации доступа использованы небуферизованные каналы, поэтому с
	собственно хранилищем в каждый момент времени работает не более одной
	процедуры, все операции с данными строго последовательные, конкурентный
	доступ исключён.

*/

import (
	"flag"
	"fmt"
	"net/http"

	"gitlab.tq-nest.lan/lancet/kvcache/vault"
	"gitlab.tq-nest.lan/lancet/kvcache/web"
)

func main() {
	ttlF := flag.Uint64("ttl", 30, "time to live, seconds")
	portF := flag.Uint("port", 8881, "port to listen")
	hostF := flag.String("host", "0.0.0.0", "specific address to listen")
	flag.Parse()
	listenStr := fmt.Sprintf("%s:%d", *hostF, *portF)
	router := new(web.Svc)
	storage := vault.Store{TTL: *ttlF}
	routing := []web.Route{
		{URL: "/storage/{id}", Methods: []string{"GET"}, Handler: router.GetValue},
		{URL: "/storage/{id}/{value}", Methods: []string{"PUT", "POST"}, Handler: router.SetValue},
		{URL: "/storage/{id}", Methods: []string{"DELETE"}, Handler: router.DelValue},
	}
	router.InitRouter(routing, &storage)
	http.Handle("/", router.GetRouter())
	fmt.Printf("Listening: %s\n", listenStr)
	http.ListenAndServe(listenStr, nil)
	return
}
