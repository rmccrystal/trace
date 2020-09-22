import React, {useState} from 'react';
import './scan.scss';
import {Button, Card, FormGroup, InputGroup} from "@blueprintjs/core";

// An input that accepts data from the barcode scanner and sends it to the server
export default function Scan() {
    let [submitted, setSubmitted] = useState(false);
    let [locationName, setLocationName] = useState('Library');

    const submit = () => {
        setSubmitted(true);
        setTimeout(() => {
            setSubmitted(false)
        }, 3000);
    }

    return <div className="scan">
        <Card className="scan-card" elevation={0}>
            {submitted
                ? <Submitted locationName={"Library"} studentName={"Ryan"} loginTime={new Date()}/>
                : <>
                    <h1 className="bp3-heading">Please scan badge to sign into {locationName}</h1>
                    <div className="bp3-text-large bp3-text-muted">If you do not have a badge, contact the current proctor</div>
                    <FormGroup
                        label="Badge ID"
                        helperText="After you scan your badge, this form will submit automatically">
                        <InputGroup large placeholder="scan" leftIcon={"align-justify"}
                                    rightElement={
                                        <Button intent={"success"} rightIcon={"arrow-right"}
                                                onClick={submit}>Submit</Button>
                                    }/>
                    </FormGroup>
                </>
            }
        </Card>
    </div>
}

// Displayed when a student scans
function Submitted({locationName, loginTime, studentName}: { locationName: string, studentName: string, loginTime: Date }) {
    return <h1 className="bp3-heading"><b>{studentName}</b> is signing out of {locationName} at {loginTime.toTimeString()}</h1>
}