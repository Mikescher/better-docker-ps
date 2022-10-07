package term

func enableColor() bool {
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return false
	}

	var mode uint32
	err = windows.GetConsoleMode(handle, &mode)
	if err != nil {
		return false
	}

	if mode&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING {
		mode = mode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		err = windows.SetConsoleMode(handle, mode)
		if err != nil {
			return false
		}
	}

	return true
}
