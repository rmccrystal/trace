import React, {useCallback, useEffect, useState} from "react";
import {getStudentsAtLocation, logoutAll, logoutStudent, TraceStudent, TraceEvent, TraceLocation} from "../api";
import {formatAMPM, onCatch, onCatchPrefix} from "./util";
import {Button, Card, HTMLTable, ICardProps, Spinner} from "@blueprintjs/core";
import moment from "moment";

export default function CurrentlyInLocation({location, ...props}: { location: TraceLocation } & ICardProps) {
    let [loading, setLoading] = useState(true);
    let [students, setStudents] = useState<{ student: TraceStudent, time: Date }[]>([]);

    const updateStudents = useCallback(() => {
        getStudentsAtLocation(location.id)
            .then(st => {
                setStudents(st);
                setLoading(false);
                setLogoutAllLoading(false);
            })
            .catch(onCatchPrefix(`Error getting student list: `));
    }, [location])

    useEffect(() => {
        setLoading(true);
        updateStudents()
    }, [location, updateStudents]);

    // TODO: Use websockets or something for this?
    useEffect(() => {
        if (!location) {
            return
        }
        let intervalID = setInterval(updateStudents, 1000);

        return () => clearInterval(intervalID)
    }, [location, updateStudents])

    const [logoutAllLoading, setLogoutAllLoading] = useState(false);
    const handleLogoutAllPress = () => {
        // This is set to false in the updateStudents func
        setLogoutAllLoading(true);
        logoutAll(location.id)
            .then(updateStudents)
            .catch(onCatch)
    }

    return <Card {...props} className="max-w-3xl w-full m-8 p-8">
        <h1 className="bp3-heading text-center">
            Currently in {location.name} ({students.length})
        </h1>
        <Button minimal className="mx-auto block my-3" onClick={handleLogoutAllPress} loading={logoutAllLoading}>
            <h4 className="bp3-text-muted bp3-heading text-center m-auto">
                Log out all
            </h4>
        </Button>
        {loading ? <Spinner className="m-8"/>
            : <Card className="p-0" elevation={1}>
                <HTMLTable condensed striped className="w-full">
                    <thead>
                    <tr>
                        <th>Name</th>
                        <th>Time in</th>
                        <th>Time Elapsed</th>
                    </tr>
                    </thead>
                    <tbody>
                    {students.sort((a, b) => a.time > b.time ? 1 : -1)
                        .map(student => <StudentRow key={student.student.id} location={location!} student={student}
                                                    onDeleteStudent={updateStudents}/>)}
                    </tbody>
                </HTMLTable>
            </Card>}
    </Card>
}

interface StudentRowProps {
    location: TraceLocation,
    student: { time: Date, student: TraceStudent },
    onDeleteStudent: () => void
}

function StudentRow({location, student, onDeleteStudent}: StudentRowProps) {
    let [logOutLoading, setLogOutLoading] = useState(false);

    return <tr>
        <td style={{verticalAlign: "middle"}}>{student.student.name}</td>
        <td style={{verticalAlign: "middle"}}>{formatAMPM(student.time)}</td>
        <td style={{verticalAlign: "middle"}}>{moment(student.time).fromNow(true)}</td>
        <td><Button style={{verticalAlign: "middle"}} className="float-right min-h-0 p-1" loading={logOutLoading} minimal text={`Log out`}
                    onClick={() => {
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
