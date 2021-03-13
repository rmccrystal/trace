import React, {useEffect, useState} from 'react';
import {ContactReport, generateContactReport, TraceLocation} from "../api";
import {Button, Card, ControlGroup, FormGroup, H4, HTMLTable, ICardProps, InputGroup, Spinner} from "@blueprintjs/core";
import {useParams, useHistory} from "react-router-dom"
import {onCatch} from "./util";

export default function ContactTracing({...props}: ICardProps) {
    let {id} = useParams();

    let inner;
    if (id) {
        inner = <ContactList id={id}/>
    } else {
        inner = <ContactTracingPrompt/>
    }

    return <Card className="max-w-3xl w-full m-8 p-8">
        {inner}
    </Card>
}

function ContactTracingPrompt() {
    let [id, setId] = useState('');
    let history = useHistory();

    let submitClick = () => {
        if (id) {
            history.push(`trace/${id}`);
        }
    }

    return <><p>View all other students that a target student has been in the same location with.</p>
        <FormGroup label="Student ID">
            <ControlGroup fill={true}>
                <InputGroup placeholder="Enter student ID" value={id} onChange={(e: any) => setId(e.target.value)}/>
                <Button icon="arrow-right" onClick={submitClick}/>
            </ControlGroup>
        </FormGroup></>
}

function ContactList({id: student_id}: { id: string }) {
    let [loading, setLoading] = useState(true);

    let [_report, setReport] = useState<ContactReport | null>(null)

    useEffect(() => {
        let end = new Date();

        // 14 days ago
        let start = new Date();
        start.setDate(end.getDate() - 14);

        generateContactReport(student_id, start, end)
            .then(setReport)
            .catch(onCatch)
            .finally(() => setLoading(false))
    }, [])

    if (loading) {
        return <Spinner/>
    }

    const report = _report!;

    return <>
        <H4>Contact report
            for <b>{report.target_student.name}</b> from {report.start_date.toString()} to {report.end_date.toString()}
        </H4>
        <HTMLTable condensed bordered interactive className="w-full">
            <thead>
            <tr>
                <th>Student</th>
                <th>Minutes in contact with {report.target_student.name}</th>
            </tr>
            </thead>
            <tbody>
            {report.contacts.map(({student, seconds_together}) => <tr>
                <td>{student.name}</td>
                <td>{seconds_together}</td>
            </tr>)}
            </tbody>
        </HTMLTable>
    </>
}