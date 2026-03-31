package match

// CheckWinner returns "X" or "O" if a winner exists.
// If there is no winner yet, it returns an empty string.
func CheckWinner(board [3][3]string) string {
	// Check rows
	for i := 0; i < 3; i++ {
		if board[i][0] != "" &&
			board[i][0] == board[i][1] &&
			board[i][1] == board[i][2] {
			return board[i][0]
		}
	}

	// Check columns
	for i := 0; i < 3; i++ {
		if board[0][i] != "" &&
			board[0][i] == board[1][i] &&
			board[1][i] == board[2][i] {
			return board[0][i]
		}
	}

	// Check main diagonal
	if board[0][0] != "" &&
		board[0][0] == board[1][1] &&
		board[1][1] == board[2][2] {
		return board[0][0]
	}

	// Check anti-diagonal
	if board[0][2] != "" &&
		board[0][2] == board[1][1] &&
		board[1][1] == board[2][0] {
		return board[0][2]
	}

	return ""
}
