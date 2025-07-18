package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}

// такая похожая задача (что интерфейс не равен нулю) уже была и
// здесь переменная err ссылается на функцию которая возвращает nil структуру,
// переменная является ссылочным типом и не равна нулю
