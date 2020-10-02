package trace

import (
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"trace/pkg/database"
)

type ContactReport struct {
	TargetStudent *database.Student
	// Contacts is an array of maps of students to the time that they have been in the same location
	// as the target student. Each element of this array represents the number of students in between that contact.
	// For example, the 1st element of this array would be the people who have been directly in contact with
	// TargetStudent, the 2nd element of this array would be the time spent with the first contact students, and so on
	Contacts []map[primitive.ObjectID]time.Duration
}

// GenerateContactReport generates a contact report for the targetStudent between startTime and endTime
func GenerateContactReport(targetStudent *database.Student, startTime time.Time, endTime time.Time, maxDepth int) (*ContactReport, error) {
	if maxDepth < 1 {
		return nil, errors.New("maxDepth must greater than 0")
	}

	events, err := database.DB.GetAllEventsBetween(startTime, endTime)
	if err != nil {
		return nil, err
	}
	students, err := database.DB.GetStudents()
	if err != nil {
		return nil, err
	}

	report := ContactReport{
		TargetStudent: targetStudent,
	}

	// Create the contact list
	report.Contacts = make([]map[primitive.ObjectID]time.Duration, maxDepth)
	// Make each map
	for n := range report.Contacts {
		report.Contacts[n] = make(map[primitive.ObjectID]time.Duration)
	}

	// Get all the students in direct contact
	for _, student := range students {
		if student.ID == targetStudent.ID {
			continue
		}

		timeWithTarget := getContactTimeWith(events, targetStudent.ID, student.ID)
		report.Contacts[0][student.ID] = timeWithTarget
	}

	// Calculate time spent for each depth.
	// This won't run if maxDepth is 1.
	// O(n^4)
	for depth := 1; depth < maxDepth; depth++ {
		// We're now getting contact time between multiple target students who have are in [depth-1]
		var targetStudents []primitive.ObjectID
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
				if student.ID == depthTargetStudent {
					continue
				}

				timeWithTarget := getContactTimeWith(events, depthTargetStudent, student.ID)
				report.Contacts[depth][student.ID] = timeWithTarget
			}
		}
	}

	return &report, nil
}

// getContactTimeWith gets the total amount of time that student1 and student2 have been in contact.
// events must be in order from oldest to newest
func getContactTimeWith(events []database.Event, student1 primitive.ObjectID, student2 primitive.ObjectID) time.Duration {
	// The total time the two students have been in contact
	var totalTime time.Duration

	// The current location of each student. Nil if they are not in a location
	var student1Location *primitive.ObjectID
	var student2Location *primitive.ObjectID

	// The time a student has entered their location
	var student1UpdateTime time.Time
	var student2UpdateTime time.Time

	// Parse each event
	for _, event := range events {
		// Continue if the event isn't with either of the students
		if event.StudentID != student1 && event.StudentID != student2 {
			continue
		}

		// True if we should subtract the two event times and update the total time
		var shouldUpdateTotal bool

		// Update the location and updateTime
		if event.EventType == database.EventEnter {
			switch event.StudentID {
			case student1:
				student1Location = &event.LocationID
				student1UpdateTime = event.Time
			case student2:
				student2Location = &event.LocationID
				student2UpdateTime = event.Time
			}
		} else if event.EventType == database.EventLeave {
			// We should update the total if someone leaves a location that the other student is in
			if student1Location != nil && student2Location != nil {
				if *student1Location == *student2Location {
					shouldUpdateTotal = true
				}
			}
			switch event.StudentID {
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
