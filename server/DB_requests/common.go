package DB_requests

func descToString(desc bool) string {
	if desc {
		return "DESC"
	}
	return ""
}

func sinceToString(leftValue, op1, op2, rightValue, since, after string, desc bool) string {
	if since != "" {
		if desc {
			return leftValue + " " + op2 + " " + rightValue + since + after
		}
		return leftValue + " " + op1 + " " + rightValue + since + after
	}
	return ""
}

func addIfNotNull(leftValue, body string) string {
	if body != "" {
		return ", " + leftValue + "'" + body + "'"
	}
	return ""
}
