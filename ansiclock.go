package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Colors struct {
	Off bool

	On       bool
	Min5     bool
	Min10    bool
	Min15    bool
	Min20    bool
	Min25    bool
	Min30    bool
	Min_to   bool
	Min_past bool
	H1       bool
	H2       bool
	H3       bool
	H4       bool
	H5       bool
	H6       bool
	H7       bool
	H8       bool
	H9       bool
	H10      bool
	H11      bool
	H12      bool
	Oclock   bool

	Ms0 bool
	Ms1 bool
	Ms2 bool
	Ms3 bool
	Ms4 bool
}

func applyFormat(base string, cols map[string]bool, col_on, col_off, reset string) string {

	_ = cols
	_ = col_on
	_ = col_off
	_ = reset

	final := strings.ReplaceAll(base, "${reset}", reset)
	for k, v := range cols {
		sub := fmt.Sprintf("${%s}", k)
		elt := col_off
		if v {
			elt = col_on
		}
		final = strings.ReplaceAll(final, sub, elt)
	}

	return final
}

func step(doClear, conky bool) {
	t := time.Now().Local()

	h := t.Hour() % 12
	m := t.Minute()

	cols := map[string]bool{
		"off": false,

		"on":       true,
		"min5":     false,
		"min10":    false,
		"min15":    false,
		"min20":    false,
		"min25":    false,
		"min30":    false,
		"min_to":   false,
		"min_past": false,
		"h1":       false,
		"h2":       false,
		"h3":       false,
		"h4":       false,
		"h5":       false,
		"h6":       false,
		"h7":       false,
		"h8":       false,
		"h9":       false,
		"h10":      false,
		"h11":      false,
		"h12":      false,
		"oclock":   false,

		"ms0": false,
		"ms1": false,
		"ms2": false,
		"ms3": false,
		"ms4": false,
	}

	if m < 5 {
		cols["oclock"] = true
	} else if m < 10 || m >= 55 {
		cols["min5"] = true
	} else if m < 15 || m >= 50 {
		cols["min10"] = true
	} else if m < 20 || m >= 45 {
		cols["min15"] = true
	} else if m < 25 || m >= 40 {
		cols["min20"] = true
	} else if m < 30 || m >= 35 {
		cols["min25"] = true
		cols["min20"] = true
		cols["min5"] = true
	} else {
		cols["min30"] = true
	}
	cols["min_to"] = m >= 35
	cols["min_past"] = !cols["min_to"] && !cols["oclock"]

	if m >= 35 {
		// form: XX to YY -> show one hour more
		h = (h + 1) % 12
	}
	cols["h1"] = h == 1
	cols["h2"] = h == 2
	cols["h3"] = h == 3
	cols["h4"] = h == 4
	cols["h5"] = h == 5
	cols["h6"] = h == 6
	cols["h7"] = h == 7
	cols["h8"] = h == 8
	cols["h9"] = h == 9
	cols["h10"] = h == 10
	cols["h11"] = h == 11
	cols["h12"] = h == 0

	ms := m % 5
	cols["ms0"] = true
	cols["ms1"] = ms > 0
	cols["ms2"] = ms > 1
	cols["ms3"] = ms > 2
	cols["ms4"] = ms > 3

	var color_on, color_off, color_reset string

	if conky {
		color_off = "{color0}"
		color_on = "{color1}"
		color_reset = ""

	} else {
		color_off = "\x1b[1;30;40m"
		color_on = "\x1b[1;37;44m"
		color_reset = "\x1b[0m"
	}

	template := `${on}I T${off} L ${on}I S${off} A S T I M E${reset}
${min15}A${off} C ${min15}Q U A R T E R${off} D C${reset}
${min20}T W E N T Y${min25} ${min5}F I V E${off} X${reset}
${min30}H A L F${off} B ${min10}T E N${off} F ${min_to}T O${reset}
${min_past}P A S T${off} E R U ${h9}N I N E${reset}
${h1}O N E${off} ${h6}S I X${off} ${h3}T H R E E${reset}
${h4}F O U R${off} ${h5}F I V E${off} ${h2}T W O${reset}
${h8}E I G H T${off} ${h11}E L E V E N${reset}
${h7}S E V E N${off} ${h12}T W E L V E${reset}
${h10}T E N${off} S E ${oclock}O C L O C K${reset}
${ms0}* ${ms1}* ${ms2}* ${ms3}* ${ms4}*${reset}`

	text := applyFormat(template, cols, color_on, color_off, color_reset)

	if doClear {
		fmt.Printf("%s", "\x1b[2J\x1b[H")
	}
	fmt.Println(text)
	fmt.Println(t.Format(time.RFC3339))
}

func showInfo() {

}

func main() {
	var conky bool
	var cont bool
	var info bool

	flag.BoolVar(&conky, "c", false, "Set if run inside conky")
	flag.BoolVar(&cont, "t", false, "Set to update automatically (clears screen)")
	flag.BoolVar(&info, "d", false, "Just description")

	flag.Parse()

	if info {
		showInfo()
		return
	}

	if cont {
		fmt.Printf("\x1b[?1049h")
		defer func() {
			fmt.Printf("\x1b[?1049l")
		}()
	}

	step(cont, conky)

	if !cont {
		return
	}

	ctx := context.Background()
	ctx, cancelFn := signal.NotifyContext(ctx, os.Interrupt)
	defer cancelFn()

	t := time.NewTicker(10 * time.Second)

mainloop:
	for {
		select {
		case <-ctx.Done():
			break mainloop
		case <-t.C:
			step(cont, conky)
		}
	}
}
