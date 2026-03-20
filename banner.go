package main

import "fmt"

// boldFg renders text in bold with a 24-bit RGB foreground color.
func boldFg(r, g, b int, text string) string {
	if noColor {
		return text
	}
	return fmt.Sprintf("\033[1;38;2;%d;%d;%dm%s\033[0m", r, g, b, text)
}

func printBanner(cfg *Config) {
	// ANSI Shadow font — 6 rows tall, letters joined with a single space gap.
	art := []string{
		`███╗   ███╗ █████╗ ███╗   ██╗████████╗██████╗  █████╗ `,
		`████╗ ████║██╔══██╗████╗  ██║╚══██╔══╝██╔══██╗██╔══██╗`,
		`██╔████╔██║███████║██╔██╗ ██║   ██║   ██████╔╝███████║`,
		`██║╚██╔╝██║██╔══██║██║╚██╗██║   ██║   ██╔══██╗██╔══██║`,
		`██║ ╚═╝ ██║██║  ██║██║ ╚████║   ██║   ██║  ██║██║  ██║`,
		`╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝`,
	}

	// Gradient: bright cyan (top) → blue → purple (bottom)
	gradient := [6][3]int{
		{0, 240, 255},
		{0, 190, 255},
		{40, 120, 255},
		{90, 60, 255},
		{150, 20, 240},
		{190, 0, 210},
	}

	fmt.Println()
	for i, line := range art {
		rgb := gradient[i]
		fmt.Println("  " + boldFg(rgb[0], rgb[1], rgb[2], line))
	}
	fmt.Println()

	// Cyan badge
	badgeText := "  ✦  dotfiles manager  ✦  manage · sync · backup  ✦  "
	if noColor {
		fmt.Println("  " + badgeText)
	} else {
		fmt.Printf("  \033[1;48;2;0;210;220m\033[38;2;0;20;40m%s\033[0m\n", badgeText)
	}
	fmt.Println()

	// Config paths
	paths := cfg.GitDir + "  →  " + cfg.WorkTree
	fmt.Println("  " + dim(paths))
	fmt.Println()

	// Hints
	fmt.Println("  " + dim("↹ tab") + "  " + dim("·") + "  " + cyan("help") + dim(" for commands") + "  " + dim("·") + "  " + dim("q / ESC to exit"))
	fmt.Println()
}
