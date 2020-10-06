import React, {useCallback, useEffect, useState} from "react";
import {
    getLocationVisits,
    getStudentsAtLocation,
    LocationVisit,
    logoutAll,
    logoutStudent,
    TraceLocation,
    TraceStudent
} from "../api";
import {Button, Card, HTMLTable, ICardProps, Spinner} from "@blueprintjs/core";
import {formatAMPM, onCatch, onCatchPrefix} from "./util";
import moment from "moment";

export default function VisitedLocationToday({location, ...props}: { location: TraceLocation } & ICardProps) {
    let [loading, setLoading] = useState(true);
    let [visits, setVisits] = useState<LocationVisit[]>([]);

    const updateStudents = useCallback(() => {
        getLocationVisits(location.id)
            .then(st => {
                setVisits(st);
                setLoading(false);
            })
            .catch(onCatchPrefix(`Error getting student list: `));
    }, [location])

    useEffect(() => {
        setLoading(true);
        updateStudents()
    }, [location, updateStudents]);

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

    return <Card {...props} className="max-w-3xl w-full m-8 p-8">
        <h1 className="bp3-heading text-center mb-6">
            Visited {location.name} today ({visits.length})
        </h1>
        <Card className="p-0" elevation={1}>
            <HTMLTable condensed striped className="w-full">
                <thead>
                <tr>
                    <th>Name</th>
                    <th>Left</th>
                    <th>Left At</th>
                </tr>
                </thead>
                <tbody>
                {visits.sort((a, b) => a.leave_time > b.leave_time ? 1 : -1)
                    .map(visit => <StudentRow key={visit.student.id} visit={visit}/>)}
                </tbody>
            </HTMLTable>
        </Card>
    </Card>
}

function StudentRow({visit}: { visit: LocationVisit }) {
    return <tr>
        <td>{visit.student.name}</td>
        <td>{moment(visit.leave_time).fromNow()}</td>
        <td>{formatAMPM(visit.leave_time)}</td>
    </tr>
}

