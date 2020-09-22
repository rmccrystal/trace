import React, {useState} from 'react';
import './scan.scss';
import {Card, FormGroup, InputGroup} from "@blueprintjs/core";

// An input that accepts data from the barcode scanner and sends it to the server
export default function Scan() {
    let [submitted, setSubmitted] = useState(true);

    const submit = () => {
        setSubmitted(true);
        setTimeout(() => {
            setSubmitted(false)
        }, 5000);
    }

    return <div className="scan">
        {submitted
            ? <Submitted locationName={"Library"} studentName={"Ryan"} loginTime={new Date()}/>
            : <Card className="scan-card" elevation={3}>
                <h1 className="bp3-heading">Library</h1>
                <FormGroup
                    label="Student Handle">
                    <InputGroup large placeholder="scan"/>
                </FormGroup>
            </Card>
        }
    </div>
}

// Displayed when a student scans
function Submitted({locationName, loginTime, studentName}: { locationName: string, studentName: string, loginTime: Date }) {
    return <Card className="submitted-card" elevation={3}>
        <h1 className="bp3-heading"><b>{ studentName }</b></h1>
    </Card>
}