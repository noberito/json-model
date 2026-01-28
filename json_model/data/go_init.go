var initialized bool

func check_model_init() {
	if !initialized {
		defer func() {
			if r := recover(); r != nil {
				panic(fmt.Sprintf("cannot initialize model checker: %v", r))
			}
		}()

CODE_BLOCK
		initialized = true
	}
}
