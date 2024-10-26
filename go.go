package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

func Calc(expression string) (float64, error) {
	postfix, err := toPostfix(expression)
	if err != nil {
		return 0, err
	}
	return evalPostfix(postfix)
}

func toPostfix(expression string) ([]string, error) {
	var postfix []string
	var stack []rune
	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])

		if unicode.IsDigit(char) {
			j := i
			for j < len(expression) && (unicode.IsDigit(rune(expression[j])) || expression[j] == '.') {
				j++
			}
			postfix = append(postfix, expression[i:j])
			i = j - 1
		} else if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				postfix = append(postfix, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("unmatched parentheses")
			}
			stack = stack[:len(stack)-1]
		} else if char == '+' || char == '-' || char == '*' || char == '/' {

			for len(stack) > 0 && precedence[stack[len(stack)-1]] >= precedence[char] {
				postfix = append(postfix, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, char)
		} else if !unicode.IsSpace(char) {
			return nil, errors.New("invalid character in expression")
		}
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, errors.New("unmatched parentheses")
		}
		postfix = append(postfix, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return postfix, nil
}

func evalPostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("invalid expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var res float64
			switch token {
			case "+":
				res = a + b
			case "-":
				res = a - b
			case "*":
				res = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("division by zero")
				}
				res = a / b
			default:
				return 0, errors.New("unknown operator")
			}
			stack = append(stack, res)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}
	return stack[0], nil
}

func main() {
	var str string
	fmt.Scan(&str)
	result, err := Calc(str)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
