package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/app/onboarding"
	"github.com/perekoshik/zoo/internal/app/reports"
	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/domain/ports"
)

type Router struct {
	cfg        config.Config
	onboard    *onboarding.Service
	animals    ports.AnimalRepository
	things     ports.ThingRepository
	invQuery   ports.InventoryQuery
	foodSvc    *reports.FoodService
	interSvc   *reports.InteractiveService
	balanceSvc *reports.BalanceService
}

func NewRouter(
	cfg config.Config,
	onboard *onboarding.Service,
	animals ports.AnimalRepository,
	things ports.ThingRepository,
	invQuery ports.InventoryQuery,
	food *reports.FoodService,
	interactive *reports.InteractiveService,
	balance *reports.BalanceService,
) *Router {
	return &Router{
		cfg:        cfg,
		onboard:    onboard,
		animals:    animals,
		things:     things,
		invQuery:   invQuery,
		foodSvc:    food,
		interSvc:   interactive,
		balanceSvc: balance,
	}
}

func (r *Router) Run(ctx context.Context) error {
	fmt.Println("Zoo CLI — commands: add thing | add animal | report balance | report food | report interactive | help | exit")
	sc := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !sc.Scan() {
			return sc.Err()
		}
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		args := fieldsKeepQuoted(line)
		cmd := strings.ToLower(args[0])

		switch cmd {
		case "exit", "quit":
			fmt.Println("bye!")
			return nil
		case "help":
			printHelp()
		case "add":
			if len(args) < 2 {
				fmt.Println("usage: add thing ... | add animal ...")
				continue
			}
			switch strings.ToLower(args[1]) {
			case "thing":
				r.cmdAddThing(ctx, args[2:])
			case "animal":
				r.cmdAddAnimal(ctx, args[2:])
			default:
				fmt.Println("unknown add subcommand:", args[1])
			}
		case "report":
			if len(args) < 2 {
				fmt.Println("usage: report balance | report food | report interactive")
				continue
			}
			switch strings.ToLower(args[1]) {
			case "balance":
				r.cmdReportBalance(ctx)
			case "food":
				r.cmdReportFood(ctx)
			case "interactive":
				r.cmdReportInteractive(ctx)
			default:
				fmt.Println("unknown report:", args[1])
			}
		default:
			fmt.Println("unknown command; type 'help'")
		}
	}
}

func (r *Router) cmdReportBalance(ctx context.Context) {
	rows, err := r.balanceSvc.List(ctx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(rows) == 0 {
		fmt.Println("(empty balance)")
		return
	}
	for _, it := range rows {
		fmt.Printf("#%d %s\n", it.InventoryNumber, it.DisplayName)
	}
}

func (r *Router) cmdReportFood(ctx context.Context) {
	sum, err := r.foodSvc.Summary(ctx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(sum.Rows) == 0 {
		fmt.Println("(no animals)")
		return
	}
	for _, row := range sum.Rows {
		fmt.Printf("#%d %-30s %3d kg/day\n", row.InventoryNumber, row.DisplayName, row.KilogramsPerDay)
	}
	fmt.Printf("TOTAL: %d kg/day\n", sum.TotalKgPerDay)
}

func (r *Router) cmdReportInteractive(ctx context.Context) {
	rows, err := r.interSvc.List(ctx)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if len(rows) == 0 {
		fmt.Println("(no animals eligible for the interactive zoo)")
		return
	}
	for _, row := range rows {
		fmt.Printf("#%d %-30s kindness %d\n", row.InventoryNumber, row.DisplayName, row.Kindness)
	}
}

func (r *Router) cmdAddThing(ctx context.Context, args []string) {
	if len(args) < 2 {
		fmt.Println("usage: add thing <type> <name>")
		fmt.Println("types: table, computer, chair, cage, feeder, waterbowl, lamp, broom, camera, trashbin, signboard")
		return
	}
	typ := strings.ToLower(args[0])
	name := strings.TrimSpace(strings.Join(args[1:], " "))
	if name == "" {
		fmt.Println("name must not be empty")
		return
	}

	var factoryName string
	switch typ {
	case "table", "computer", "chair", "cage", "feeder", "waterbowl", "lamp", "broom", "camera", "trashbin", "signboard":
		factoryName = strings.Title(typ) + " " + name
	default:
		fmt.Println("unknown thing type:", typ)
		return
	}

	wrap, err := newGenericThing(factoryName)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	stored, err := r.things.Add(ctx, wrap)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("OK: added thing #%d %s\n", stored.InventoryNumber(), stored.DisplayName())
}

func (r *Router) cmdAddAnimal(ctx context.Context, args []string) {
	if len(args) < 3 {
		fmt.Println("usage: add animal <species> <name> <foodKgPerDay> [kindness-for-herbivore]")
		return
	}
	species := canonicalSpecies(args[0])
	name := args[1]
	food, err := strconv.Atoi(args[2])
	if err != nil || food < 0 {
		fmt.Println("foodKgPerDay must be non-negative integer")
		return
	}

	isHerb := isHerbivoreSpecies(species)
	var entity ports.AnimalEntity

	if isHerb {
		if len(args) < 4 {
			fmt.Println("herbivore requires [kindness 0..10]")
			return
		}
		kindness, err := strconv.Atoi(args[3])
		if err != nil || kindness < 0 || kindness > 10 {
			fmt.Println("kindness must be 0..10")
			return
		}
		h, err := animal.NewHerbivore(name, species, food, kindness)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		entity = h
	} else {
		p, err := animal.NewPredator(name, species, food)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		entity = p
	}

	res, err := r.onboard.Admit(ctx, entity)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	if !res.Accepted {
		fmt.Printf("DENIED: %s — reason: %s\n", entity.DisplayName(), res.Reason)
		return
	}
	fmt.Printf("OK: admitted #%d %s (food %d kg/day)\n",
		res.Animal.InventoryNumber(), res.Animal.DisplayName(), res.Animal.FoodPerDayKg())
}

func printHelp() {
	fmt.Println(`commands:
  add thing <type> <name>              - add an inventory thing
  add animal <species> <name> <food> [kindness] - add an animal via vet clinic
  report balance                       - list all items on balance
  report food						   - list all food requirements
  report interactive				   - list all animal kindness
  help | exit`)
}

func fieldsKeepQuoted(s string) []string {
	var out []string
	var cur strings.Builder
	inQuote := false
	for _, r := range s {
		switch r {
		case '"':
			inQuote = !inQuote
		case ' ':
			if inQuote {
				cur.WriteRune(r)
			} else if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}

func canonicalSpecies(s string) string {
	return strings.Title(strings.ToLower(strings.TrimSpace(s)))
}

func isHerbivoreSpecies(species string) bool {
	switch strings.ToLower(species) {
	case "rabbit", "monkey", "elephant", "giraffe", "zebra",
		"panda", "kangaroo", "flamingo", "parrot", "deer",
		"hippopotamus", "rhinoceros":
		return true
	default:
		return false
	}
}

type genericThing struct {
	name string
	num  int
}

func newGenericThing(name string) (*genericThing, error) {
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name must not be empty")
	}
	return &genericThing{name: name}, nil
}

func (t *genericThing) Name() string         { return t.name }
func (t *genericThing) DisplayName() string  { return t.name }
func (t *genericThing) InventoryNumber() int { return t.num }
func (t *genericThing) AssignInventoryNumber(n int) error {
	if n < 1 {
		return fmt.Errorf("inventory number must be >= 1")
	}
	if t.num != 0 {
		return fmt.Errorf("inventory number already set (%d)", t.num)
	}
	t.num = n
	return nil
}
