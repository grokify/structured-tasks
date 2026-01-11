package main

import (
	"fmt"
	"sort"

	"github.com/grokify/structured-roadmap/roadmap"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats <file>",
	Short: "Show roadmap statistics",
	Long:  `Display statistics about items, statuses, and categories in a roadmap.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runStats,
}

func runStats(cmd *cobra.Command, args []string) error {
	path := args[0]

	r, err := roadmap.ParseFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	stats := r.Stats()
	out := cmd.OutOrStdout()

	fmt.Fprintf(out, "Roadmap: %s\n", r.Project)
	fmt.Fprintf(out, "Total items: %d\n\n", stats.Total)

	// Status breakdown
	fmt.Fprintln(out, "By Status:")
	statusOrder := []roadmap.Status{roadmap.StatusCompleted, roadmap.StatusInProgress, roadmap.StatusPlanned, roadmap.StatusFuture}
	legend := r.GetLegend()
	for _, status := range statusOrder {
		count := stats.ByStatus[status]
		if count > 0 {
			pct := float64(count) / float64(stats.Total) * 100
			entry := legend[status]
			fmt.Fprintf(out, "  %s %s: %d (%.0f%%)\n", entry.Emoji, entry.Description, count, pct)
		}
	}

	// Priority breakdown
	if len(stats.ByPriority) > 0 {
		fmt.Fprintln(out, "\nBy Priority:")
		priorityOrder := []roadmap.Priority{roadmap.PriorityCritical, roadmap.PriorityHigh, roadmap.PriorityMedium, roadmap.PriorityLow}
		for _, priority := range priorityOrder {
			count := stats.ByPriority[priority]
			if count > 0 {
				pct := float64(count) / float64(stats.Total) * 100
				fmt.Fprintf(out, "  %s: %d (%.0f%%)\n", roadmap.PriorityLabel(priority), count, pct)
			}
		}
	}

	// Area breakdown
	if len(stats.ByArea) > 0 {
		fmt.Fprintln(out, "\nBy Area:")
		// Sort areas by count
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
			for _, area := range r.Areas {
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
		// Sort types by count
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
	if len(r.Phases) > 0 {
		fmt.Fprintln(out, "\nBy Phase:")
		itemsByPhase := r.ItemsByPhase()
		for _, phase := range r.Phases {
			items := itemsByPhase[phase.ID]
			status := ""
			if phase.Status != "" {
				status = " " + r.GetStatusEmoji(phase.Status)
			}
			fmt.Fprintf(out, "  %s%s: %d items\n", phase.Name, status, len(items))
		}
	}

	// Progress
	fmt.Fprintf(out, "\nProgress: %.0f%% complete\n", stats.CompletedPercent())
	return nil
}
