package trace

import (
	"errors"
	"github.com/sirupsen/logrus"
	"time"
	"trace/pkg/database"
)

type ContactReport struct {
	TargetStudent *database.Student `json:"target_student"`
	// Contacts is an array of maps of students to the time that they have been in the same location
	// as the target student. Each element of this array represents the number of students in between that contact.
	// For example, the 1st element of this array would be the people who have been directly in contact with
	// TargetStudent, the 2nd element of this array would be the time spent with the first contact students, and so on
	Contacts []map[database.StudentRef]time.Duration `json:"contacts"`
}

// GenerateContactReport generates a contact report for the targetStudent between startTime and endTime
func GenerateContactReport(targetStudent *database.Student, startTime time.Time, endTime time.Time, maxDepth int) (*ContactReport, error) {
	if maxDepth < 1 {
		return nil, errors.New("maxDepth must greater than 0")
	}

	events := database.DB.GetAllEventsBetween(startTime, endTime)
	students := database.DB.GetStudents()

	report := ContactReport{
		TargetStudent: targetStudent,
	}

	// Create the contact list
	report.Contacts = make([]map[database.StudentRef]time.Duration, maxDepth)
	// Make each map
	for n := range report.Contacts {
		report.Contacts[n] = make(map[database.StudentRef]time.Duration)
	}

	// Get all the students in direct contact
	for _, student := range students {
		if student.ID == targetStudent.ID {
			continue
		}

		timeWithTarget := getContactTimeWith(events, targetStudent.Ref(), student.Ref())
		if timeWithTarget == 0 {
			continue
		}
		report.Contacts[0][student.Ref()] = timeWithTarget
	}

	// Calculate time spent for each depth.
	// This won't run if maxDepth is 1.
	// O(n^4)
	for depth := 1; depth < maxDepth; depth++ {
		// We're now getting contact time between multiple target students who have are in [depth-1]
		var targetStudents []database.StudentRef
		for student := range report.Contacts[depth-1] {
			targetStudents = append(targetStudents, student)
		}

		// If there are no more targetStudents, stop
		if len(targetStudents) == 0 {
			break
		}

		// For all students who have been in contact with target with the targetStudent [depth-1] away
		for _, depthTargetStudent := range targetStudents {
			for _, student := range students {
				if student.Ref() == depthTargetStudent {
					continue
				}

				timeWithTarget := getContactTimeWith(events, depthTargetStudent, student.Ref())
				report.Contacts[depth][student.Ref()] = timeWithTarget
			}
		}
	}

	return &report, nil
}

// getContactTimeWith gets the total amount of time that student1 and student2 have been in contact.
// events must be in order from oldest to newest
func getContactTimeWith(events []database.Event, student1 database.StudentRef, student2 database.StudentRef) time.Duration {
	// The total time the two students have been in contact
	var totalTime time.Duration

	// The current location of each student. Nil if they are not in a location
	var student1Location *database.LocationRef
	var student2Location *database.LocationRef

	// The time a student has entered their location
	var student1UpdateTime time.Time
	var student2UpdateTime time.Time

	// Parse each event
	for _, event := range events {
		// Continue if the event isn't with either of the students
		if event.Student != student1 && event.Student != student2 {
			continue
		}

		// True if we should subtract the two event times and update the total time
		var shouldUpdateTotal bool

		// Update the location and updateTime
		if event.EventType == database.EventEnter {
			switch event.Student {
			case student1:
				student1Location = &event.Location
				student1UpdateTime = event.Time
			case student2:
				student2Location = &event.Location
				student2UpdateTime = event.Time
			}
		} else if event.EventType == database.EventLeave {
			// We should update the total if someone leaves a location that the other student is in
			if student1Location != nil && student2Location != nil {
				if *student1Location == *student2Location {
					shouldUpdateTotal = true
				}
			}
			switch event.Student {
			case student1:
				student1Location = nil
				student1UpdateTime = event.Time
			case student2:
				student2Location = nil
				student2UpdateTime = event.Time
			}
		} else {
			logrus.Warnf("Encountered invalid event type %d in getContactTimeWith", event.EventType)
			continue
		}
		// TODO: Add location timeouts

		if shouldUpdateTotal {
			// Get the time spent together in the location
			timeSpentTogether := student1UpdateTime.Sub(student2UpdateTime)

			// Take absolute value of timeSpentTogether
			if timeSpentTogether < 0 {
				timeSpentTogether = -timeSpentTogether
			}

			totalTime += timeSpentTogether
		}
	}

	return totalTime
}
