package trace

// HandleScan should be called whenever a student scans in or scans out.
// It will return the Events that it creates or an error.
// If the studentID cannot be found in the database, it will not be stored
//func HandleScan(studentID string) (database.Events, error) {
//
//}