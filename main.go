package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Node struct {
	token string
	lhs   *Node
	rhs   *Node
}

func main() {
	vars := map[string]float64{}
	input := get_input()
	process_input(input, vars)
	for k, v := range vars {
		fmt.Printf("%v = %v\n", k, v)
	}
}

// Change this to get the appropriate input (via stdin and such)
func get_input() string {
	return "x = 10 + 3 * 5\ny = 50 + x * 30\nz = 6 * 3 + 2 / 4"
}

// Process the input string containing the \n-delimited statements
func process_input(input string, vars map[string]float64) {
	statements := strings.Split(input, "\n")
	for _, s := range statements {
		process_statement(s, vars)
	}
}

// process a statement of the form: y = 1 + 3 * 5 + 4 / 2 + 4 * x
func process_statement(s string, vars map[string]float64) {
	node := string_to_tree(s)
	print("postfix: ")
	print_postfix(node)
	print("\n")
	process_recursive(node, vars)
}

// Conversion of the string to a binary tree with Op, lhs & rhs
func string_to_tree(s string) *Node {
	tokens := strings.Split(s, " ")
	nodes := make([]*Node, len(tokens))
	n_tokens := len(tokens) // Number of tokens in the equation
	for i, t := range tokens {
		// For each token, create a node
		nodes[i] = new(Node)
		*nodes[i] = Node{t, nil, nil}
	}
	ops := []string{"^", "/", "*", "-", "+", "="}
	for _, c := range ops {
		// Assuming that the string starts with either a number or a variable
		for i := 0; i < n_tokens; i++ {
			if nodes[i] == nil || nodes[i].token != c {
				continue // skip everything if node at index i is nil
			}
			j, k := i-1, i+1 //select previous and next elements
			for j >= 0 && nodes[j] == nil {
				j--
			}
			for k < n_tokens && nodes[k] == nil {
				k++
			}
			// Build the tree and replace pointers in the array
			nodes[i].lhs = nodes[j]
			nodes[j] = nil
			nodes[i].rhs = nodes[k]
			nodes[k] = nil
		}
	}
	// At this point, the nodes slice is guaranteed
	// to have only one non-nil element in it.
	i := 0
	for nodes[i] == nil {
		i++
	}
	node := nodes[i]
	return node
}

// Recursively compute the tree representing the equation
func process_recursive(node *Node, vars map[string]float64) float64 {
	switch node.token {
	case "=": // Assign the computed rhs (operation or value) to the lhs (variable)
		vars[node.lhs.token] = process_recursive(node.rhs, vars)
		return 0
	case "+": // Addition
		return process_recursive(node.lhs, vars) + process_recursive(node.rhs, vars)
	case "-": // Subtraction
		return process_recursive(node.lhs, vars) - process_recursive(node.rhs, vars)
	case "*": // Multiplication
		return process_recursive(node.lhs, vars) * process_recursive(node.rhs, vars)
	case "/": // Division
		return process_recursive(node.lhs, vars) / process_recursive(node.rhs, vars)
	case "^": // Exponent
		return math.Pow(process_recursive(node.lhs, vars), process_recursive(node.rhs, vars))
	default:
		// Parse the node's number value if required
		val, err := strconv.ParseFloat(node.token, 64)
		if err != nil { // If node isn't a number
			var_val, pres := vars[node.token] // Get variable value
			if !pres {
				fmt.Println("error parsing: " + node.token)
				return 0
			}
			val = var_val
		}
		return val
	}
}

// Debugging utilities (print equation in postfix notation)
func print_postfix(root *Node) {
	if root.lhs != nil {
		print_postfix(root.lhs)
	}
	if root.rhs != nil {
		print_postfix(root.rhs)
	}
	print(root.token + " ")
}
