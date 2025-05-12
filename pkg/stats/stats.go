package stats

import (
	"fmt"
	"sort"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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
	// Normalize both dates to the beginning of their respective days
	date = GetBeginningOfDay(date)
	now := GetBeginningOfDay(time.Now())

	// Calculate the difference in days
	diff := now.Sub(date)
	days := int(diff.Hours() / HoursInDay)

	if days > DaysInLastSixMonths {
		return OutOfRange
	}
	return days
}

// CalculateWeekdayOffset calculates an offset value based on the current day of the week.
// This is used for positioning in the contribution graph.
//
// Returns:
//   - int: A value from 0 to 6 representing the day of the week (0=Sunday, 1=Monday, etc.)
func CalculateWeekdayOffset() int {
	weekday := time.Now().Weekday()

	switch weekday {
	case time.Sunday:
		return 0
	case time.Monday:
		return 1
	case time.Tuesday:
		return 2
	case time.Wednesday:
		return 3
	case time.Thursday:
		return 4
	case time.Friday:
		return 5
	case time.Saturday:
		return 6
	}

	return 0 // Should never reach here
}

// GetCommitsFromRepo retrieves commit information from a Git repository.
// If an email is provided, it filters commits by that email address.
// If no email is provided, it includes commits from all users.
// It updates the provided commits map with the count of commits per day.
//
// Parameters:
//   - email: The email address to filter commits by (if empty, includes all commits)
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

	// Get the commit history starting from HEAD
	iterator, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		return nil, fmt.Errorf("failed to get commit log: %w", err)
	}

	// Iterate through the commits
	err = iterator.ForEach(func(c *object.Commit) error {
		// If email is provided, skip commits not authored by the specified email
		if email != "" && c.Author.Email != email {
			return nil
		}

		daysAgo := CountDaysSinceDate(c.Author.When)

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

// ProcessRepositories processes a Git repository and collects commit statistics.
// If an email is provided, it filters commits by that email address.
// If no email is provided, it includes commits from all users.
//
// Parameters:
//   - email: The email address to filter commits by (if empty, includes all commits)
//   - directory: The directory to analyze (should be a Git repository)
//
// Returns:
//   - map[int]int: A map of days to commit counts
//   - error: An error if any occurred during processing
func ProcessRepositories(email string, directory string) (map[int]int, error) {
	// Initialize the commits' map with zeros for all days
	commits := make(map[int]int, DaysInLastSixMonths)
	for i := DaysInLastSixMonths; i > 0; i-- {
		commits[i] = 0
	}

	// Process the repository
	var err error
	commits, err = GetCommitsFromRepo(email, directory, commits)
	if err != nil {
		return nil, fmt.Errorf("error processing repository at %s: %w", directory, err)
	}

	return commits, nil
}

// PrintCell prints a single cell in the contribution graph with the appropriate coloring
// based on the number of commits and whether it represents today.
//
// Parameters:
//   - val: The number of commits for this cell
//   - today: Whether this cell represents today
//   - date: The date for this cell
//   - showCommitCount: Whether to display the number of commits on each cell
//   - showDaysOfMonth: Whether to display the days of the month on the graph calendar
func PrintCell(val int, today bool, date time.Time, showCommitCount bool, showDaysOfMonth bool) {
	// Light gray for no contributions
	escape := "\033[0;37;48;5;248m"

	// Set color based on commit count - from lighter to darker green
	switch {
	case val > 0 && val < 5:
		escape = "\033[1;30;48;5;120m" // Light green for few commits
	case val >= 5 && val < 10:
		escape = "\033[1;30;48;5;34m" // Medium green for moderate commits
	case val >= 10:
		escape = "\033[1;30;48;5;22m" // Dark green for many commits
	}

	// Special color for today's cell
	if today {
		escape = "\033[1;37;45m"
	}

	// Determine what to display in the cell
	cellContent := "   " // Default empty cell

	// Show the commit count if requested
	if showCommitCount && val > 0 {
		if val < 10 {
			cellContent = fmt.Sprintf(" %d ", val) // Single digit with padding
		} else {
			cellContent = fmt.Sprintf("%d ", val) // Double-digit with padding
		}
	}

	// Show day of the month if requested
	if showDaysOfMonth {
		day := date.Day()
		if day < 10 {
			cellContent = fmt.Sprintf(" %d ", day) // Single digit with padding
		} else {
			cellContent = fmt.Sprintf("%d ", day) // Double-digit with padding
		}
	}

	// Print cell with a pipe separator
	fmt.Printf("%s%s%s|", escape, cellContent, "\033[0m")
}

// PrintCommitsStats displays a visual representation of commit statistics in a calendar-like grid.
// It processes the commits' map, builds the columns, and prints the cells.
//
// Parameters:
//   - commits: A map of days to commit counts
//   - showCommitCount: Whether to display the number of commits on each cell
//   - showDaysOfMonth: Whether to display the days of the month on the graph calendar
func PrintCommitsStats(commits map[int]int, showCommitCount bool, showDaysOfMonth bool) {
	keys := SortMapIntoSlice(commits)
	cols := BuildCols(keys, commits)
	PrintCells(cols, showCommitCount, showDaysOfMonth)
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

	// Get today's date
	today := GetBeginningOfDay(time.Now())

	// Initialize a map to group commits by week and day
	weekDayCommits := make(map[int]map[int]int)

	// Calculate the current weekday
	_ = int(today.Weekday())

	// Calculate the start date for the contribution graph (6 months ago)
	startDate := today.AddDate(0, -6, 0)

	// Calculate the start of the week for the start date
	daysToStartSunday := int(startDate.Weekday())
	startOfFirstWeek := startDate.AddDate(0, 0, -daysToStartSunday)

	for _, k := range keys {
		// Calculate the actual date for this key (days ago)
		date := today.AddDate(0, 0, -k)

		// Skip dates before the start date
		if date.Before(startDate) {
			continue
		}

		// Get the actual day of the week for this date (0=Sunday, 1=Monday, etc.)
		dayInWeek := int(date.Weekday())

		// Calculate the number of weeks since the start of the first week
		weeksSinceStart := int(date.Sub(startOfFirstWeek).Hours() / (HoursInDay * DaysInWeek))

		// The week number is the number of weeks from the start of the graph
		week := WeeksInLastSixMonths - weeksSinceStart

		// Initialize the week map if it doesn't exist
		if _, ok := weekDayCommits[week]; !ok {
			weekDayCommits[week] = make(map[int]int)
		}

		// Add the commit count to the week/day map
		weekDayCommits[week][dayInWeek] += commits[k]
	}

	// Convert the week/day map to columns
	for week, days := range weekDayCommits {
		col := make(Column, 7) // Initialize with 7 days (0-6)

		// Fill in the column with commit counts for each day
		for day, count := range days {
			col[day] = count
		}

		cols[week] = col
	}

	return cols
}

// calculateGraphParameters calculates the parameters needed for rendering the contribution graph.
// It determines which week today is in and the maximum week to display.
//
// Parameters:
//   - cols: A map of week numbers to columns of commit counts
//
// Returns:
//   - time.Time: The start of the first week in the graph
//   - int: The week number that contains today
//   - int: The maximum week number to display
func calculateGraphParameters(cols map[int]Column) (time.Time, int, int) {
	// Calculate which week today is in
	today := GetBeginningOfDay(time.Now())
	startDate := today.AddDate(0, -6, 0)
	daysToStartSunday := int(startDate.Weekday())
	startOfFirstWeek := startDate.AddDate(0, 0, -daysToStartSunday)
	weeksSinceStart := int(today.Sub(startOfFirstWeek).Hours() / (HoursInDay * DaysInWeek))
	todayWeek := WeeksInLastSixMonths - weeksSinceStart

	// Find the maximum week number in the col map
	maxWeek := 0
	for week := range cols {
		if week > maxWeek {
			maxWeek = week
		}
	}

	// Ensure we display at least WeeksInLastSixMonths+1 columns
	if maxWeek < WeeksInLastSixMonths {
		maxWeek = WeeksInLastSixMonths
	}

	return startOfFirstWeek, todayWeek, maxWeek
}

// printCellForPosition prints the appropriate cell for a given position in the contribution graph.
//
// Parameters:
//   - cols: A map of week numbers to columns of commit counts
//   - weekNum: The week number for this cell
//   - dayNum: The day number for this cell
//   - todayWeek: The week number that contains today
//   - cellDate: The date for this cell
//   - showCommitCount: Whether to display the number of commits on each cell
//   - showDaysOfMonth: Whether to display the days of the month on the graph calendar
func printCellForPosition(cols map[int]Column, weekNum int, dayNum int, todayWeek int, cellDate time.Time, showCommitCount bool, showDaysOfMonth bool) {
	// Check if this cell represents today
	isToday := weekNum == todayWeek && dayNum == CalculateWeekdayOffset()

	// Get a commit count for this cell if available
	commitCount := 0
	if col, ok := cols[weekNum]; ok && len(col) > dayNum {
		commitCount = col[dayNum]
	}

	// Print the cell with appropriate styling
	PrintCell(commitCount, isToday, cellDate, showCommitCount, showDaysOfMonth)
}

// printWeekRow prints a single row (day of the week) in the contribution graph.
//
// Parameters:
//   - cols: A map of week numbers to columns of commit counts
//   - dayNum: The day number (0-6) to print
//   - startOfFirstWeek: The start date of the first week in the graph
//   - todayWeek: The week number that contains today
//   - maxWeek: The maximum week number to display
//   - showCommitCount: Whether to display the number of commits on each cell
//   - showDaysOfMonth: Whether to display the days of the month on the graph calendar
func printWeekRow(cols map[int]Column, dayNum int, startOfFirstWeek time.Time, todayWeek int, maxWeek int, showCommitCount bool, showDaysOfMonth bool) {
	// Iterate through weeks (columns)
	for weekNum := maxWeek + 1; weekNum >= 0; weekNum-- {
		// Print day labels in the first column
		if weekNum == maxWeek+1 {
			PrintDayCol(dayNum)
			continue
		}

		// Calculate the date for this cell
		weekOffset := WeeksInLastSixMonths - weekNum
		cellDate := startOfFirstWeek.AddDate(0, 0, weekOffset*7+dayNum)

		// Print the appropriate cell for this position
		printCellForPosition(cols, weekNum, dayNum, todayWeek, cellDate, showCommitCount, showDaysOfMonth)
	}
	fmt.Printf("\n")
}

// PrintCells renders the contribution graph by printing all cells in a grid format.
// It first prints the month labels, then iterates through each day of the week and each week,
// printing the appropriate cell for each position.
//
// Parameters:
//   - cols: A map of week numbers to columns of commit counts
//   - showCommitCount: Whether to display the number of commits on each cell
//   - showDaysOfMonth: Whether to display the days of the month on the graph calendar
func PrintCells(cols map[int]Column, showCommitCount bool, showDaysOfMonth bool) {
	PrintMonths()

	// Calculate graph parameters
	startOfFirstWeek, todayWeek, maxWeek := calculateGraphParameters(cols)

	// Iterate through days of the week (rows)
	for dayNum := 0; dayNum <= 6; dayNum++ {
		printWeekRow(cols, dayNum, startOfFirstWeek, todayWeek, maxWeek, showCommitCount, showDaysOfMonth)
	}
}

// PrintMonths prints the month labels at the top of the contribution graph.
// It places month names on columns with the first day of that month.
func PrintMonths() {
	// Started from 6 months ago
	startDate := GetBeginningOfDay(time.Now()).AddDate(0, -6, 0)

	// Calculate the start of the week for the start date
	daysToSunday := int(startDate.Weekday())
	startOfWeek := startDate.AddDate(0, 0, -daysToSunday)

	// Print initial spacing
	fmt.Printf("         ")

	// Map to store week numbers that contain the first day of a month
	monthLabels := make(map[int]string)

	// Iterate through each day in the 6-month period to find the first days of months
	for weekNum := WeeksInLastSixMonths; weekNum >= 0; weekNum-- {
		for dayInWeek := 0; dayInWeek < 7; dayInWeek++ {
			// Calculate the date for this cell
			cellDate := startOfWeek.AddDate(0, 0, (WeeksInLastSixMonths-weekNum)*7+dayInWeek)

			// If this is the first day of a month, store the month label for the previous week
			if cellDate.Day() == 1 {
				// Only store the label if we're not at the oldest week (to avoid out of bounds)
				if weekNum < WeeksInLastSixMonths {
					monthLabels[weekNum+1] = cellDate.Month().String()[:3]
				}
				break // Found first day of the month in this week, move to next week
			}
		}
	}

	// Print month labels
	for weekNum := WeeksInLastSixMonths; weekNum >= 0; weekNum-- {
		if label, ok := monthLabels[weekNum]; ok {
			fmt.Printf("%s ", label)
		} else {
			fmt.Printf("    ")
		}
	}

	// Add an extra column for the current week
	fmt.Printf("    ")

	fmt.Printf("\n")
}

// PrintDayCol prints the day labels on the left side of the contribution graph.
// It displays the first letter of each day of the week.
//
// Parameters:
//   - day: The day index (0-6) to print a label for
func PrintDayCol(day int) {
	out := "     " // Default empty label

	switch day {
	case 0:
		out = "  S  "
	case 1:
		out = "  M  "
	case 2:
		out = "  T  "
	case 3:
		out = "  W  "
	case 4:
		out = "  T  "
	case 5:
		out = "  F  "
	case 6:
		out = "  S  "
	}

	fmt.Printf("%s", out)
}
