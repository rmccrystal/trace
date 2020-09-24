import React, {useEffect, useState} from "react";
import {Button, Card, HTMLTable, Spinner} from "@blueprintjs/core";
import {getStudentsAtLocation, Student, TraceEvent} from "../api";
import {useGlobalState} from "../app";
import {formatAMPM, onCatch} from "./util";
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
        getStudentsAtLocation(location.id)
            .then(st => {
                setLoading(false);
                setStudents(st);
            })
            .catch(onCatch)
    }, [location]);

    // TODO: Use websockets or something for this?
    useEffect(() => {
        if (!location) {
            return
        }
        let intervalID = setInterval(() => {
            if (!location) {
                return
            }
            getStudentsAtLocation(location.id)
                .then(st => {
                    setStudents(st);
                })
                .catch(onCatch)
        }, 1000);

        return () => clearInterval(intervalID)
    }, [location])

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
            {/*<th>Log Out</th>*/}
            </thead>
            <tbody>
            {students.sort((a, b) => a.event.time > b.event.time ? 1 : -1).map(student => <tr>
                <td style={{verticalAlign: "middle"}}>{student.student.name}</td>
                <td style={{verticalAlign: "middle"}}>{formatAMPM(student.event.time)}</td>
                <td style={{verticalAlign: "middle"}}>{moment(student.event.time).fromNow(true)}</td>
                <td><Button icon="delete" className="float-right" minimal text={`Log out`}/></td>
            </tr>)}
            </tbody>
        </HTMLTable>
    </Card>
}