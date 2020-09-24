import React, {useState} from 'react';
import './scan.scss';
import {Button, Card, FormGroup, InputGroup} from "@blueprintjs/core";
import {EventType, scan, TraceEvent} from "../api";
import {onCatch} from "./util";

// An input that accepts data from the barcode scanner and sends it to the server
export default function Scan() {
    let [locationName, setLocationName] = useState('Library');
    let [state, setState] = useState<"form" | "submitted" | "loading">("form");
    let [event, setEvent] = useState<TraceEvent | null>(null);

    let [handle, setHandle] = useState("");
    const handleChange = (e: any) => {
        setHandle(e.target.value);
    }

    const submit = () => {
        setState("loading");
        scan(handle, "5f658215319b49a857b35831")
            .then((ev) => {
                setEvent(ev);
                setState("submitted");
                setTimeout(() => {
                    setState("form")
                }, 3000);
            })
            .catch((e: any) => {
                onCatch(e);
                setState("form");
            })
    }

    const handleKeyDown = (e: any) => {
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
                <InputGroup large onChange={handleChange} onKeyDown={handleKeyDown} placeholder=""
                            leftIcon={"align-justify"}
                            rightElement={<Button minimal rightIcon={"arrow-right"} loading={state === "loading"}
                                                  onClick={submit}/>}/>
            </FormGroup>
        </>;
    } else if (state === "submitted") {
        contentElem = <Submitted event={event!}/>
    }


    return <div className="flex items-center justify-center w-full h-full">
        <Card className="scan-card" elevation={0}>
            {contentElem}
        </Card>
    </div>
}

// Displayed when a student scans
function Submitted({event}: { event: TraceEvent }) {
    const {event_type, location_name, student_name, time} = event;
    return <div>
        <h1 className="bp3-heading text-5xl">Hello <b>{student_name}</b>! You are currently
            checking <b>{event_type === EventType.Enter ? "in to " : "out of "}</b>
            the <b>{location_name}</b>.</h1> <br/>
        <p className="bp3-text-large bp3-text-muted text-xl">If this is not you, please see the proctor on duty</p>
    </div>
}