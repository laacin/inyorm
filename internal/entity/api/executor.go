package api

type Executor[Prep any] interface {
	Run(binder ...Scanner) error
	Prepare(fn func(exec Prep) error) error
}

type Prepare interface {
	Run(args []Value, scanner ...Scanner) error
}
