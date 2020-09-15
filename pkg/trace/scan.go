package trace

import "trace/pkg/database"

// HandleScan should be called whenever a student scans in or scans out.
// It will return the Event that it creates or an error.
// If the studentID cannot be found in the database, it will NOT be stored
func HandleScan(studentID string) (database.Event, error) {

}