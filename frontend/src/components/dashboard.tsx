import React, {useEffect, useState} from "react";
import {Button, Card, HTMLTable, Spinner} from "@blueprintjs/core";
import {getStudentsAtLocation, Student} from "../api";
import {useGlobalState} from "../app";
import {onCatch} from "./util";

// todo: preserve state while changing the page back
export default function Dashboard() {
    let [loading, setLoading] = useState(true);
    let [students, setStudents] = useState<Student[]>([]);

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

    if(loading) {
        return <Spinner className="mt-10"/>
    }

    return <Card className="max-w-xl w-full m-8">
        <h1 className="bp3-heading text-center">Currently in {location?.name} ({students.length})</h1>
        <HTMLTable condensed striped className="w-full">
            <thead>
                <th>Name</th>
                <th>Time Elapsed</th>
                {/*<th>Log Out</th>*/}
            </thead>
            <tbody>
            {students.map(student => <tr>
                <td style={{verticalAlign: "middle"}}>{student.name}</td>
                <td style={{verticalAlign: "middle"}}>1 hour</td>
                <td><Button icon="delete" className="float-right" minimal text={`Log out`}/></td>
            </tr>)}
            </tbody>
        </HTMLTable>
    </Card>
}