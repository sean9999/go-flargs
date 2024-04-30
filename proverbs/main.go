package proverbs

import (
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/sean9999/go-flargs"
)

//go:embed proverbs.txt
var Proverbs string

const NumberOfProverbs = 19

type Params struct {
	index     int
	remaining []string
}

func (p *Params) Args() []string {
	return p.remaining
}

func (p *Params) Parse(args []string) error {
	if len(args) < 1 {
		p.remaining = args
	} else {
		i, err := strconv.Atoi(args[0])
		if err != nil {
			p.remaining = args
			return err
		} else {
			p.index = i
			p.remaining = args[1:]
		}
	}
	return nil
}
func (p *Params) Load(env *flargs.Environment) error {
	if p.index == 0 {
		r := rand.New(env.Randomness)
		p.index = r.Intn(NumberOfProverbs)
	}
	if NumberOfProverbs <= p.index {
		return errors.New("out of range")
	}
	return nil
}
func (p *Params) Run(env *flargs.Environment) error {
	proverb := strings.Split(strings.TrimSpace(Proverbs), "\n")[p.index]
	proverb = fmt.Sprintln(proverb)
	_, err := env.OutputStream.Write([]byte(proverb))
	return err
}
