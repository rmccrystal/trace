import React, {useCallback, useEffect, useState} from "react";
import {getLocationVisits, LocationVisit, TraceLocation} from "../api";
import {Card, HTMLTable, ICardProps, Spinner} from "@blueprintjs/core";
import {formatAMPM, onCatchPrefix} from "./util";
import moment from "moment";
import {Simulate} from "react-dom/test-utils";

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

    return <Card {...props} className="max-w-3xl w-full m-8 p-8">
        <h1 className="bp3-heading text-center mb-6">
            Visited {location.name} Today ({visits.length})
        </h1>
        {loading
            ? <Spinner className="mt-8"/>
            : <VisitTable visits={visits}/>
        }
    </Card>
}

function VisitTable({visits}: { visits: LocationVisit[] }) {
    return <Card className="p-0" elevation={1}>
        <HTMLTable condensed striped className="w-full">
            <thead>
            <tr>
                <th>Name</th>
                <th>Left</th>
                <th>Time In</th>
                <th>Time Out</th>
                <th>Time Spent Last Visit</th>
            </tr>
            </thead>
            <tbody>
            {visits.sort((a, b) => a.leave_time < b.leave_time ? 1 : -1)
                .map(visit => <VisitRow key={visit.student.id} visit={visit}/>)}
            </tbody>
        </HTMLTable>
    </Card>
}

function VisitRow({visit}: { visit: LocationVisit }) {
    return <tr>
        <td>{visit.student.name}</td>
        <td>{moment(visit.leave_time).fromNow()}</td>
        <td>{formatAMPM(visit.enter_time)}</td>
        <td>{formatAMPM(visit.leave_time)}</td>
        <td>{moment(visit.enter_time).from(visit.leave_time, true)}</td>
    </tr>
}

