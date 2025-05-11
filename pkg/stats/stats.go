package stats

import (
	"fmt"
	"sort"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"git-contrib/pkg/fileutil"
)

// Constants for time calculations and display
const (
	OutOfRange           = 99999
	DaysInLastSixMonths  = 183
	WeeksInLastSixMonths = 26
	HoursInDay           = 24
	DaysInWeek           = 7
)

type Column []int

// GetBeginningOfDay returns a new time.Time with the same date as the input time
// but with the time set to 00:00:00.
//
// Parameters:
//   - t: The time to get the beginning of the day for
//
// Returns:
//   - time.Time: A new time.Time representing the beginning of the day
func GetBeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// CountDaysSinceDate calculates the number of days between the given date and today.
// If the difference is greater than DaysInLastSixMonths, it returns OutOfRange.
//
// Parameters:
//   - date: The starting date to count from
//
// Returns:
//   - int: The number of days since the given date, or OutOfRange if more than DaysInLastSixMonths
func CountDaysSinceDate(date time.Time) int {
	days := 0
	now := GetBeginningOfDay(time.Now())

	for date.Before(now) {
		date = date.Add(time.Hour * HoursInDay)
		days++
		if days > DaysInLastSixMonths {
			return OutOfRange
		}
	}
	return days
}

// CalculateWeekdayOffset calculates an offset value based on the current day of the week.
// This is used for positioning in the contribution graph.
//
// Returns:
//   - int: A value from 1 to 7 representing the offset for the current weekday
func CalculateWeekdayOffset() int {
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		return 7
	case time.Monday:
		return 6
	case time.Tuesday:
		return 5
	case time.Wednesday:
		return 4
	case time.Thursday:
		return 3
	case time.Friday:
		return 2
	case time.Saturday:
		return 1
	}

	return 0 // Should never reach here
}

// GetCommitsFromRepo retrieves commit information from a Git repository for a specific email address.
// It updates the provided commits map with the count of commits per day.
//
// Parameters:
//   - email: The email address to filter commits by
//   - path: The path to the Git repository
//   - commits: A map of days to commit counts to update
//
// Returns:
//   - map[int]int: The updated commits map
//   - error: An error if any occurred during repository processing
func GetCommitsFromRepo(email string, path string, commits map[int]int) (map[int]int, error) {
	// Open the git repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository at %s: %w", path, err)
	}

	// Get the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Get the commits history starting from HEAD
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	// Calculate offset for the contribution graph
	offset := CalculateWeekdayOffset()

	// Iterate through the commits
	err = iterator.ForEach(func(c *object.Commit) error {
		// Skip commits not authored by the specified email
		if c.Author.Email != email {
			return nil
		}

		daysAgo := CountDaysSinceDate(c.Author.When) + offset

		// Only count commits within the last six months
		if daysAgo != OutOfRange {
			commits[daysAgo]++
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error processing commits: %w", err)
	}

	return commits, nil
}

// ProcessRepositories processes all repositories listed in the .git-contrib dotfile
// and collects commit statistics for the specified email address.
//
// Parameters:
//   - email: The email address to filter commits by
//
// Returns:
//   - map[int]int: A map of days to commit counts
//   - error: An error if any occurred during processing
func ProcessRepositories(email string) (map[int]int, error) {
	filePath := fileutil.GetDotfilePath()
	repos := fileutil.ParseFileLines(filePath)

	// Initialize the commits' map with zeros for all days
	commits := make(map[int]int, DaysInLastSixMonths)
	for i := DaysInLastSixMonths; i > 0; i-- {
		commits[i] = 0
	}

	// Process each repository
	for _, path := range repos {
		var err error
		commits, err = GetCommitsFromRepo(email, path, commits)
		if err != nil {
			return nil, fmt.Errorf("error processing repository %s: %w", path, err)
		}
	}

	return commits, nil
}

// PrintCell prints a single cell in the contribution graph with appropriate coloring
// based on the number of commits and whether it represents today.
//
// Parameters:
//   - val: The number of commits for this cell
//   - today: Whether this cell represents today
func PrintCell(val int, today bool) {
	// Default color for empty cells
	escape := "\033[0;37;30m"

	// Set color based on commit count
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;47m" // Light color for few commits
	case val >= 5 && val < 10:
		escape = "\033[1;30;43m" // Medium color for moderate commits
	case val >= 10:
		escape = "\033[1;30;42m" // Dark color for many commits
	}

	// Special color for today's cell
	if today {
		escape = "\033[1;37;45m"
	}

	// Print empty cell
	if val == 0 {
		fmt.Printf(escape + "  - " + "\033[0m")
		return
	}

	// Format string based on number of digits
	str := "  %d "
	switch {
	case val >= 10 && val < 100:
		str = " %d "
	case val >= 100:
		str = "%d "
	}

	fmt.Printf(escape+str+"\033[0m", val)
}

// PrintCommitsStats displays a visual representation of commit statistics in a calendar-like grid.
// It processes the commits map, builds the columns, and prints the cells.
//
// Parameters:
//   - commits: A map of days to commit counts
func PrintCommitsStats(commits map[int]int) {
	keys := SortMapIntoSlice(commits)
	cols := BuildCols(keys, commits)
	PrintCells(cols)
}

// SortMapIntoSlice extracts the keys from a map and returns them as a sorted slice.
//
// Parameters:
//   - m: The map to extract keys from
//
// Returns:
//   - []int: A sorted slice of the map's keys
func SortMapIntoSlice(m map[int]int) []int {
	var keys []int
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

// BuildCols organizes commit data into columns for display in the contribution graph.
// Each column represents a week, and each cell in the column represents a day.
//
// Parameters:
//   - keys: A sorted slice of day indices
//   - commits: A map of days to commit counts
//
// Returns:
//   - map[int]Column: A map of week numbers to columns of commit counts
func BuildCols(keys []int, commits map[int]int) map[int]Column {
	cols := make(map[int]Column)
	col := Column{}

	for _, k := range keys {
		// Calculate week number and day of week
		week := int(k / DaysInWeek) // 26, 25, ..., 1
		dayInWeek := k % DaysInWeek // 0, 1, 2, 3, 4, 5, 6

		// Start a new column at the beginning of each week
		if dayInWeek == 0 {
			col = Column{}
		}

		// Add commit count to the column
		col = append(col, commits[k])

		// Save the column when we reach the end of the week
		if dayInWeek == 6 {
			cols[week] = col
		}
	}

	return cols
}

// PrintCells renders the contribution graph by printing all cells in a grid format.
// It first prints the month labels, then iterates through each day of the week and each week,
// printing the appropriate cell for each position.
//
// Parameters:
//   - cols: A map of week numbers to columns of commit counts
func PrintCells(cols map[int]Column) {
	PrintMonths()

	// Iterate through days of the week (rows)
	for j := 6; j >= 0; j-- {
		// Iterate through weeks (columns)
		for i := WeeksInLastSixMonths + 1; i >= 0; i-- {
			// Print day labels in the first column
			if i == WeeksInLastSixMonths+1 {
				PrintDayCol(j)
			}

			if col, ok := cols[i]; ok {
				// Special case for today's cell
				if i == 0 && j == CalculateWeekdayOffset()-1 {
					PrintCell(col[j], true)
					continue
				} else if len(col) > j {
					PrintCell(col[j], false)
					continue
				}
			}
			// Print empty cell if no data
			PrintCell(0, false)
		}
		fmt.Printf("\n")
	}
}

// PrintMonths prints the month labels at the top of the contribution graph.
// It calculates the appropriate position for each month label based on the current date.
func PrintMonths() {
	// Start from 6 months ago
	week := GetBeginningOfDay(time.Now()).Add(-(DaysInLastSixMonths * time.Hour * HoursInDay))
	month := week.Month()

	// Print initial spacing
	fmt.Printf("         ")

	// Print month names when month changes
	for {
		if week.Month() != month {
			fmt.Printf("%s ", week.Month().String()[:3])
			month = week.Month()
		} else {
			fmt.Printf("    ")
		}

		// Move to next week
		week = week.Add(DaysInWeek * time.Hour * HoursInDay)
		if week.After(time.Now()) {
			break
		}
	}
	fmt.Printf("\n")
}

// PrintDayCol prints the day labels on the left side of the contribution graph.
// It displays labels for Monday, Wednesday, and Friday.
//
// Parameters:
//   - day: The day index (0-6) to print a label for
func PrintDayCol(day int) {
	out := "     " // Default empty label

	switch day {
	case 1:
		out = " Mon "
	case 3:
		out = " Wed "
	case 5:
		out = " Fri "
	}

	fmt.Printf(out)
}
