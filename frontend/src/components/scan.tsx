import React, {useState} from 'react';
import './scan.scss';
import {Button, Card, FormGroup, InputGroup, Spinner} from "@blueprintjs/core";

// An input that accepts data from the barcode scanner and sends it to the server
export default function Scan() {
    let [locationName, setLocationName] = useState('Library');
    let [state, setState] = useState<"form" | "submitted" | "loading">("form");

    const submit = () => {
        setState("loading");
        setTimeout(() => {
            setState("submitted");
            setTimeout(() => {
                setState("form")
            }, 3000);
        }, 1000)
    }

    const onBadgeInputKeyDown = (e: any) => {
        if (e.key === "Enter") {
            submit();
        }
    }

    // The element inside the container card
    let contentElem;
    if (state === "form" || state === "loading") {
        contentElem = <>
            <h1 className="bp3-heading">Please scan badge to sign into {locationName}</h1>
            <div className="bp3-text-large bp3-text-muted mb-5">If you do not have a badge, contact the current
                proctor
            </div>
            <FormGroup
                label="Badge ID"
                helperText="After you scan your badge, this form will submit automatically">
                <InputGroup large onKeyDown={onBadgeInputKeyDown} placeholder="" leftIcon={"align-justify"}
                            rightElement={<Button minimal rightIcon={"arrow-right"} loading={state === "loading"}
                                                  onClick={submit}/>}/>
            </FormGroup>
        </>;
    } else if (state === "submitted") {
        contentElem = <Submitted locationName={"Library"} studentName={"Ryan McCrystal"} loginTime={new Date()}/>
    }


    return <div className="flex items-center justify-center w-full h-full">
        <Card className="scan-card" elevation={0}>
            {contentElem}
        </Card>
    </div>
}

// Displayed when a student scans
function Submitted({locationName, loginTime, studentName}: { locationName: string, studentName: string, loginTime: Date }) {
    return <div>
        <h1 className="bp3-heading text-5xl">Hello <b>{studentName}</b>! You are currently checking out of
            the <b>{locationName}</b>.</h1> <br/>
        <p className="bp3-text-large bp3-text-muted text-xl">If this is not you, please see the proctor on duty</p>
    </div>
}