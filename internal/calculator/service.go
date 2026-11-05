package calculator

// Service - сервис калькулятора. Сейчас полей у него нет, но позже сюда
// добавятся зависимости (например, база данных для истории операций).
type Service struct{}

// NewService - конструктор: создаёт сервис и возвращает на него указатель.
func NewService() *Service {
	return &Service{}
}

// Calculate проверяет запрос и выполняет операцию над двумя числами.
func (s *Service) Calculate(req Request) (Response, error) {
	// Отрицательные суммы не принимаем (решение из шага 0).
	if req.A < 0 || req.B < 0 {
		return Response{}, ErrNegativeAmount
	}

	switch req.Op {
	case OpAdd:
		return Response{Result: req.A + req.B}, nil
	case OpSub:
		return Response{Result: req.A - req.B}, nil
	case OpMul:
		return Response{Result: req.A * req.B}, nil
	case OpDiv:
		// Проверка нужна до деления: на целых числах деление на 0 роняет программу.
		if req.B == 0 {
			return Response{}, ErrDivisionByZero
		}
		return Response{Result: req.A / req.B, Remainder: req.A % req.B}, nil
	default:
		return Response{}, ErrUnknownOperation
	}

}
