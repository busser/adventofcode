package d20

import (
	"fmt"
	"io"
	"strings"

	"github.com/busser/adventofcode/helpers"
)

// PartOne solves the first problem of day 20 of Advent of Code 2023.
func PartOne(r io.Reader, w io.Writer) error {
	modules, err := modulesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	n := newNetwork(modules)

	for i := 0; i < 1000; i++ {
		n.pressButton()
	}

	product := n.lowPulseCount * n.highPulseCount

	_, err = fmt.Fprintf(w, "%d", product)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

// PartTwo solves the second problem of day 20 of Advent of Code 2023.
func PartTwo(r io.Reader, w io.Writer) error {
	modules, err := modulesFromReader(r)
	if err != nil {
		return fmt.Errorf("could not read input: %w", err)
	}

	n := newNetwork(modules)

	count, err := n.countPressesToTurnOn()
	if err != nil {
		return fmt.Errorf("could not count presses: %w", err)
	}

	_, err = fmt.Fprintf(w, "%d", count)
	if err != nil {
		return fmt.Errorf("could not write answer: %w", err)
	}

	return nil
}

const (
	helloPulse = iota
	lowPulse
	highPulse
)

type pulse struct {
	from string
	to   string
	kind int
}

type module interface {
	receive(p pulse, send func(pulse))
	hello(send func(pulse))
	getName() string
	listDestinations() []string
}

var (
	_ module = &broadcastModule{}
	_ module = &flipFlopModule{}
	_ module = &conjuctionModule{}
)

type network struct {
	modules map[string]module

	lowPulseCount  int
	highPulseCount int

	highPulseSent map[string]bool
}

func newNetwork(modules map[string]module) *network {
	n := network{
		modules:       modules,
		highPulseSent: make(map[string]bool),
	}

	// Sending hello pulses allows conjuction modules to know their inputs.
	for _, m := range modules {
		m.hello(n.send)
	}

	return &n
}

func (n *network) send(p pulse) {
	var nextPulses []pulse
	newPulse := func(p pulse) {
		nextPulses = append(nextPulses, p)

		switch p.kind {
		case lowPulse:
			n.lowPulseCount++
		case highPulse:
			n.highPulseCount++
			n.highPulseSent[p.from] = true
		}
	}

	newPulse(p)

	for len(nextPulses) > 0 {
		p := nextPulses[0]
		nextPulses = nextPulses[1:]

		target, ok := n.modules[p.to]
		if !ok {
			// pulse sent into the void
			continue
		}

		target.receive(p, newPulse)
	}
}

func (n *network) pressButton() {
	n.send(pulse{
		from: "button",
		to:   "broadcaster",
		kind: lowPulse,
	})
}

func (n *network) moduleInputs(name string) []string {
	var inputs []string
	for _, m := range n.modules {
		for _, d := range m.listDestinations() {
			if d == name {
				inputs = append(inputs, m.getName())
			}
		}
	}
	return inputs
}

func (n *network) countPressesToTurnOn() (int, error) {
	// This function relies on the input's structure.
	// It is not a general solution.
	//
	// The machine turns on when the rx module receives a low pulse.
	// The rx module has a single input, a conjuction module.
	// That conjuction module has multiple inputs, all conjuction modules.
	//
	// The rx module's input sends a low pulse when all of its inputs send a
	// high pulse. These inputs send a high pulse cyclically.
	//
	// This solution identifies the inputs of the rx module's input, determines
	// their cycle lenght, and then calculates the least common multiple of
	// those cycle lengths. The result is the number of times the button must
	// be pressed to turn on the machine.

	rxInputs := n.moduleInputs("rx")
	if len(rxInputs) != 1 {
		return 0, fmt.Errorf("unexpected number of inputs for rx module: %d", len(rxInputs))
	}

	cyclicalInputs := n.moduleInputs(rxInputs[0])
	if len(cyclicalInputs) == 0 {
		return 0, fmt.Errorf("unexpected number of inputs for %s module: %d", rxInputs[0], len(cyclicalInputs))
	}

	cycleLengths := make(map[string]int)

	for count := 1; ; count++ {
		n.pressButton()

		for _, mod := range cyclicalInputs {
			if n.highPulseSent[mod] {
				if _, ok := cycleLengths[mod]; !ok {
					cycleLengths[mod] = count
				}
			}
		}

		allCyclesFound := len(cycleLengths) == len(cyclicalInputs)
		if allCyclesFound {
			fullCycleLength := 1
			for _, cycleLength := range cycleLengths {
				fullCycleLength = lcm(fullCycleLength, cycleLength)
			}
			return fullCycleLength, nil
		}
	}
}

func lcm(a, b int) int {
	if a == 0 && b == 0 {
		return 0
	}

	return abs(a*b) / gcd(a, b)
}

func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type baseModule struct {
	name         string
	destinations []string
}

func (m *baseModule) hello(send func(pulse)) {
	for _, d := range m.destinations {
		send(pulse{
			from: m.name,
			to:   d,
			kind: helloPulse,
		})
	}
}

func (m *baseModule) getName() string {
	return m.name
}

func (m *baseModule) listDestinations() []string {
	return m.destinations
}

type broadcastModule struct {
	baseModule
}

func newBroadcastModule(name string, destinations []string) *broadcastModule {
	return &broadcastModule{
		baseModule: baseModule{
			name:         name,
			destinations: destinations,
		},
	}
}

func (m *broadcastModule) receive(p pulse, send func(pulse)) {
	switch p.kind {
	case helloPulse:
	case lowPulse, highPulse:
		for _, d := range m.destinations {
			send(pulse{
				from: p.to,
				to:   d,
				kind: p.kind,
			})
		}
	}
}

type flipFlopModule struct {
	baseModule
	on bool
}

func newFlipFlopModule(name string, destinations []string) *flipFlopModule {
	return &flipFlopModule{
		baseModule: baseModule{
			name:         name,
			destinations: destinations,
		},
		on: false,
	}
}

func (m *flipFlopModule) receive(p pulse, send func(pulse)) {
	switch p.kind {
	case helloPulse:
	case lowPulse:
		m.on = !m.on
		nextKind := lowPulse
		if m.on {
			nextKind = highPulse
		}
		for _, d := range m.destinations {
			send(pulse{
				from: p.to,
				to:   d,
				kind: nextKind,
			})
		}
	case highPulse:
		return
	default:
		panic(fmt.Sprintf("unknown pulse kind: %d", p.kind))
	}
}

type conjuctionModule struct {
	baseModule
	recentPulses map[string]int
}

func newConjuctionModule(name string, destinations []string) *conjuctionModule {
	return &conjuctionModule{
		baseModule: baseModule{
			name:         name,
			destinations: destinations,
		},
		recentPulses: make(map[string]int),
	}
}

func (m *conjuctionModule) receive(p pulse, send func(pulse)) {
	switch p.kind {
	case helloPulse:
		m.recentPulses[p.from] = lowPulse
	case lowPulse, highPulse:
		m.recentPulses[p.from] = p.kind

		allHigh := true
		for _, p := range m.recentPulses {
			if p != highPulse {
				allHigh = false
				break
			}
		}

		nextKind := highPulse
		if allHigh {
			nextKind = lowPulse
		}
		for _, d := range m.destinations {
			send(pulse{
				from: p.to,
				to:   d,
				kind: nextKind,
			})
		}
	default:
		panic(fmt.Sprintf("unknown pulse kind: %d", p.kind))
	}
}

func modulesFromReader(r io.Reader) (map[string]module, error) {
	lines, err := helpers.LinesFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("could not read input: %w", err)
	}

	modules := make(map[string]module)
	for _, line := range lines {
		parts := strings.SplitN(line, " -> ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		if len(parts[0]) < 2 {
			return nil, fmt.Errorf("invalid line: %s", line)
		}

		destinations := strings.Split(parts[1], ", ")

		switch {
		case parts[0] == "broadcaster":
			name := parts[0]
			modules[name] = newBroadcastModule(name, destinations)
		case parts[0][0] == '%':
			name := parts[0][1:]
			modules[name] = newFlipFlopModule(name, destinations)
		case parts[0][0] == '&':
			name := parts[0][1:]
			modules[name] = newConjuctionModule(name, destinations)
		default:
			return nil, fmt.Errorf("invalid line: %s", line)
		}
	}

	return modules, nil
}
