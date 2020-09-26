import React, {useEffect, useState} from "react";
import {Button, Card, HTMLTable, Spinner} from "@blueprintjs/core";
import {getStudentsAtLocation, logoutStudent, Student, TraceEvent, TraceLocation} from "../api";
import {formatAMPM, onCatch, onCatchPrefix} from "./util";
import moment from "moment";

// todo: preserve state while changing the page back
export default function Dashboard({location}: {location: TraceLocation}) {
    let [loading, setLoading] = useState(true);
    let [students, setStudents] = useState<{ event: TraceEvent, student: Student }[]>([]);

    useEffect(() => {
        setLoading(true);
        updateStudents()
    }, [location]);

    const updateStudents = () => {
        getStudentsAtLocation(location.id)
            .then(st => {
                setStudents(st);
                setLoading(false);
            })
            .catch(onCatchPrefix(`Error getting student list`));
    }

    // List of logout buttons that should be loading
    let [loadingLogoutButtonIDs, setLoadingLogoutButtonIDs] = useState<string[]>([]);

    // TODO: Use websockets or something for this?
    useEffect(() => {
        if (!location) {
            return
        }
        let intervalID = setInterval(updateStudents, 1000);

        return () => clearInterval(intervalID)
    }, [location, updateStudents])

    if (loading) {
        return <Spinner className="mt-10"/>
    }

    return <Card className="max-w-xl w-full m-8">
        <h1 className="bp3-heading text-center">Currently in {location?.name} ({students.length})</h1>
        <HTMLTable condensed striped className="w-full">
            <thead>
            <th>Name</th>
            <th>Time in</th>
            <th>Time Elapsed</th>
            <th/>
            </thead>
            <tbody>
            {students.sort((a, b) => a.event.time > b.event.time ? 1 : -1)
                .map(student => <StudentRow key={student.student.id} location={location!} student={student} onDeleteStudent={updateStudents}/>)}
            </tbody>
        </HTMLTable>
    </Card>
}

interface StudentRowProps {
    location: TraceLocation,
    student: { event: TraceEvent, student: Student },
    onDeleteStudent: () => void
}

function StudentRow({location, student, onDeleteStudent}: StudentRowProps) {
    let [logOutLoading, setLogOutLoading] = useState(false);

    return <tr>
        <td style={{verticalAlign: "middle"}}>{student.student.name}</td>
        <td style={{verticalAlign: "middle"}}>{formatAMPM(student.event.time)}</td>
        <td style={{verticalAlign: "middle"}}>{moment(student.event.time).fromNow(true)}</td>
        <td><Button icon="delete" className="float-right" loading={logOutLoading} minimal text={`Log out`} onClick={() => {
            if (!location) {
                return
            }

            setLogOutLoading(true)
            logoutStudent(student.student.id, location!.id)
                .finally(() => {
                    onDeleteStudent();
                })
                .catch(() => {
                    onCatchPrefix(`Error logging out ${student.student.name}: `);
                    setLogOutLoading(false);
                });
        }}/></td>
    </tr>
}