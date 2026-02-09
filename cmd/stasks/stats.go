package main

import (
	"fmt"
	"sort"

	"github.com/grokify/structured-tasks/tasks"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats <file>",
	Short: "Show task list statistics",
	Long:  `Display statistics about tasks, statuses, and categories in a task list.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStats,
}

func runStats(cmd *cobra.Command, args []string) error {
	path := args[0]

	tl, err := tasks.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	stats := tl.Stats()
	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "Task List: %s\n", tl.Project)
	fmt.Fprintf(out, "Total tasks: %d\n\n", stats.Total)

	// Status breakdown
	fmt.Fprintln(out, "By Status:")
	legend := tl.GetLegend()
	for _, status := range tasks.StatusOrder() {
		count := stats.ByStatus[status]
		if count > 0 {
			pct := float64(count) / float64(stats.Total) * 100
			entry := legend[status]
			fmt.Fprintf(out, "  %s %s: %d (%.0f%%)\n", entry.Emoji, entry.Description, count, pct)
		}
	}

	// Area breakdown
	if len(stats.ByArea) > 0 {
		fmt.Fprintln(out, "\nBy Area:")
		type areaCount struct {
			name  string
			count int
		}
		areas := make([]areaCount, 0, len(stats.ByArea))
		for area, count := range stats.ByArea {
			areas = append(areas, areaCount{area, count})
		}
		sort.Slice(areas, func(i, j int) bool {
			return areas[i].count > areas[j].count
		})
		for _, a := range areas {
			// Find area name
			name := a.name
			for _, area := range tl.Areas {
				if area.ID == a.name {
					name = area.Name
					break
				}
			}
			fmt.Fprintf(out, "  %s: %d\n", name, a.count)
		}
	}

	// Type breakdown
	if len(stats.ByType) > 0 {
		fmt.Fprintln(out, "\nBy Type:")
		type typeCount struct {
			name  string
			count int
		}
		types := make([]typeCount, 0, len(stats.ByType))
		for t, count := range stats.ByType {
			types = append(types, typeCount{t, count})
		}
		sort.Slice(types, func(i, j int) bool {
			return types[i].count > types[j].count
		})
		for _, t := range types {
			fmt.Fprintf(out, "  %s: %d\n", t.name, t.count)
		}
	}

	// Phase breakdown
	phases := tl.PhaseNumbers()
	if len(phases) > 0 {
		fmt.Fprintln(out, "\nBy Phase:")
		tasksByPhase := tl.TasksByPhase()
		for _, phase := range phases {
			phaseTasks := tasksByPhase[phase]
			fmt.Fprintf(out, "  Phase %d: %d tasks\n", phase, len(phaseTasks))
		}
		// Unphased tasks
		if unphasedTasks := tasksByPhase[0]; len(unphasedTasks) > 0 {
			fmt.Fprintf(out, "  Unphased: %d tasks\n", len(unphasedTasks))
		}
	}

	// Progress
	completedPct := 0.0
	if stats.Total > 0 {
		completedPct = float64(stats.CompletedCount()) / float64(stats.Total) * 100
	}
	fmt.Fprintf(out, "\nProgress: %.0f%% complete\n", completedPct)
	return nil
}
