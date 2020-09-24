import React, {useEffect, useState} from "react";
import {Card, Spinner} from "@blueprintjs/core";
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

    return <Card className="max-w-xl m-auto mt-8">
        <h1 className="bp3-heading">Students in {location?.name}</h1>
        {students.map(student => <div className="bp3-text-large">{student.name}</div>)}
    </Card>
}