package calculator

import "errors"

// Operation - своя строка для операции. Так в коде видно, что здесь ждём именно
// операцию, а не любой текст.
type Operation string

// Операции, которые умеет калькулятор. Константы нужны, чтобы не писать "add"
// руками: опечатка в строке не заметна, а в имени константы - сразу ошибка
const (
	OpAdd Operation = "add" // сложение: A + B
	OpSub Operation = "sub" // вычитание: A - B
	OpMul Operation = "mul" // умножение: A * B
	OpDiv Operation = "div" // деление нацело: A / B, остаток идёт отдельно
)

// Ошибки калькулятора. Создаём один раз, чтобы обработчик мог узнать ошибку
// через errors.Is и ответить клиенту нужным кодом.
var (
	ErrUnknownOperation = errors.New("unknow operation") // операции нет в списке выше
	ErrDivisionByZero   = errors.New("divizion by zero") // делить на ноль нельзя
	ErrNegativeAmount   = errors.New("negative amount")  // сумма не может быть меньше нуля
)

// Request - что присылает клиент, например {"a": 100, "b": 3, "op": "div"}.
// Числа целые (int64), а не с запятой: у float бывают погрешности, а деньгам они не нужны.
type Request struct {
	A  int64     `json: "a"`  // первое число
	B  int64     `json: "b"`  // второе число
	Op Operation `json: "op:` // что с ними сделать
}

// Response - что мы отвечаем, например {"result": 33, "remainder": 1}.
type Response struct {
	Result    int64 `json: "result"`      // итог; при делении - целая часть
	Remainder int64 `json: : "remainder"` // остаток от деления, в остальных операциях 0
}
