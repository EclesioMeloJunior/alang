package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/EclesioMeloJunior/ducklang/eval"
	"github.com/EclesioMeloJunior/ducklang/lexer"
	"github.com/EclesioMeloJunior/ducklang/object"
	"github.com/EclesioMeloJunior/ducklang/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnv()

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		io.WriteString(out, "\t"+err.Error()+"\n")
	}
}
