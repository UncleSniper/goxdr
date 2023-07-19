package goxdr

type EmptyReadState struct {}

func(state EmptyReadState) Update([]byte) (int, bool) {
	return 0, true
}

func(state EmptyReadState) EndPacket() error {
	return nil
}

var _ ReadState = EmptyReadState{}

var TheEmptyReadState EmptyReadState
