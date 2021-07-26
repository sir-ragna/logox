package main

import (
	"fmt"
	"log"
)

/* p => q
 * q => p
 * p v q
 * p ^ q
 * p == q (equivalent)
 * p EQUIV p
 * (p => q) v k
 * q AND p
 * p OR q
 * p AND (p OR k) OR s
 * ~p
 * NOT p
 * p IMPL q
 * p IMPL (p OR k)
 * ~ (p AND q)
 * p NAND q
 * ~ (p OR q)
 * p NOR q
 * p v q
 * ~ (p v q)
 * ~ (p ^ q)
 */

type Token int

const (
	ILLEGAL     Token = iota
	EOF               // End Of File
	WS                // White space
	SYMBOL            // p, q, k, ...
	IMPLICATION       // =>
	OR                //   OR   v
	AND               //  AND  ^
	NOT               //  NOT  ~
	NOR               //  NOR
	NAND              // NAND
	XOR               //  XOR
	TRUE              // TRUE  1
	FALSE             // FALSE 0
)

type Node struct {
	token Token
	name  string
	nodes []Node
}

func (n *Node) evaluate(context map[string]bool) bool {
	switch n.token {
	case SYMBOL:
		val, exist := context[n.name]
		if exist {
			return val
		} else {
			log.Fatalf("SYMBOL %q cannot be found\n", n.name)
		}
	case TRUE:
		return true
	case FALSE:
		return false
	case NOT:
		if len(n.nodes) != 1 {
			log.Fatalln("Cannot negate more than one term")
		}
		return !n.nodes[0].evaluate(context)
	case IMPLICATION:
		/* Truth table: https://dyclassroom.com/boolean-algebra/propositional-logic-truth-table
		 */
		if len(n.nodes) != 2 {
			log.Fatalln("Implication can only have two terms")
		}
		// p => q
		if n.nodes[0].evaluate(context) == true { // eval p
			if n.nodes[1].evaluate(context) == true { // eval q
				// This must be true as well
				return true
			} else {
				// p is true but q is not
				// This is impossible
				fmt.Printf("Failed implication \"%s(%t) => %s(%t)\"\n", n.nodes[0].name, n.nodes[0].evaluate(context), n.nodes[1].name, n.nodes[1].evaluate(context))
				return false
			}
		} else {
			// false => ____
			return true
		}
	case OR:
		if len(n.nodes) != 2 {
			log.Fatalln("OR can only have two terms")
		}
		return n.nodes[0].evaluate(context) || n.nodes[1].evaluate(context)
	case AND:
		if len(n.nodes) != 2 {
			log.Fatalln("AND can only have two terms")
		}
		return n.nodes[0].evaluate(context) && n.nodes[1].evaluate(context)
	}

	return false
}

func main() {
	var context = make(map[string]bool)

	// The AST for "p => q"
	astExample := &Node{
		token: IMPLICATION,
		name:  "p => q",
		nodes: []Node{
			Node{token: SYMBOL, name: "p"},
			Node{token: SYMBOL, name: "q"},
		},
	}
	// figure out the list of symbols

	// list of symbols
	var symbols = []string{"p", "q"}

	for n := 0; n < (1 << len(symbols)); n++ {
		// 1 << len(symbols) means _2 to the power of len(symbols)_
		for index, symbol := range symbols {
			context[symbol] = n&(1<<index) > 0
			// index:          0  1  2  3  4
			// 1 << index:     1  2  4  8 16
			// n:              0  1  2  3  4

			// n&(1<<index)
			//-------------
			// 0 & (1 << 0) == 0      0 & (1 << 1) == 0
			// 1 & (1 << 0) == 1      1 & (1 << 1) == 0
			// 2 & (1 << 0) == 0      2 & (1 << 1) == 2
			// 3 & (1 << 0) == 1      3 & (1 << 1) == 2

			// > 0
			//----
			// 00
			// 10
			// 01
			// 11
		}

		result := astExample.evaluate(context)
		for key, val := range context {
			fmt.Printf("%s=%t ", key, val)
		}
		fmt.Printf("  \t Evaluating \"%s\" :: %t\n", astExample.name, result)
	}

}
