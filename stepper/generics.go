package stepper

func reverseSequence(slice [8][4]uint8) [8][4]uint8 {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = reverseSequenceLine(slice[j]), reverseSequenceLine(slice[i])
	}
	return slice
}

func reverseSequenceLine(line [4]uint8) [4]uint8 {
	for i, j := 0, len(line)-1; i < j; i, j = i+1, j-1 {
		line[i], line[j] = line[j], line[i]
	}
	return line
}
