package entities

type Executor interface {
	Run(binder ...Scanner) error
	Prepare(fn func(exec Prepare) error) error
}

type Prepare interface {
	Run(args []Value, scanner ...Scanner) error
}
