import React, {useEffect, useState} from "react";
import {Button, Card, HTMLTable, Spinner} from "@blueprintjs/core";
import {getStudentsAtLocation, logoutStudent, Student, TraceEvent} from "../api";
import {useGlobalState} from "../app";
import {formatAMPM, onCatch, onCatchPrefix} from "./util";
import moment from "moment";

// todo: preserve state while changing the page back
export default function Dashboard() {
    let [loading, setLoading] = useState(true);
    let [students, setStudents] = useState<{ event: TraceEvent, student: Student }[]>([]);

    let [location] = useGlobalState('location')

    useEffect(() => {
        setLoading(true);
        if (!location) {
            return
        }
        updateStudents()
    }, [location]);

    const updateStudents = () => {
        if (!location) {
            return
        }
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
    }, [location])

    if (loading || location === undefined) {
        return <Spinner className="mt-10"/>
    }


    return <Card className="max-w-xl w-full m-8">
        <h1 className="bp3-heading text-center">Currently in {location?.name} ({students.length})</h1>
        <HTMLTable condensed striped className="w-full">
            <thead>
            <th>Name</th>
            <th>Time in</th>
            <th>Time Elapsed</th>
            {/*<th>Log Out</th>*/}
            </thead>
            <tbody>
            {students.sort((a, b) => a.event.time > b.event.time ? 1 : -1).map(student => <tr>
                <td style={{verticalAlign: "middle"}}>{student.student.name}</td>
                <td style={{verticalAlign: "middle"}}>{formatAMPM(student.event.time)}</td>
                <td style={{verticalAlign: "middle"}}>{moment(student.event.time).fromNow(true)}</td>
                <td><Button icon="delete" className="float-right" minimal text={`Log out`} onClick={() => {
                    if(!location) {
                        return
                    }

                    setLoadingLogoutButtonIDs([...loadingLogoutButtonIDs, student.student.id]);
                    logoutStudent(student.student.id, location!.id)
                        .then(() => {
                            updateStudents();
                            setLoadingLogoutButtonIDs(loadingLogoutButtonIDs.filter(id => id != student.student.id));
                        })
                        .catch(onCatchPrefix(`Error logging out ${student.student.name}: `));
                }}/></td>
            </tr>)}
            </tbody>
        </HTMLTable>
    </Card>
}