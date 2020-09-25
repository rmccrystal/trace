import React, {useEffect, useRef, useState} from 'react';
import './scan.scss';
import {Button, Card, FormGroup, InputGroup, Spinner} from "@blueprintjs/core";
import {EventType, scan, TraceEvent} from "../api";
import {onCatch} from "./util";
import {useGlobalState} from '../app';

// An input that accepts data from the barcode scanner and sends it to the server
export default function Scan() {
    let [state, setState] = useState<"form" | "submitted" | "loading">("form");
    let [event, setEvent] = useState<TraceEvent | null>(null);

    let [handle, setHandle] = useState("");
    const handleChange = (e: any) => {
        setHandle(e.target.value);
    }

    let [_location] = useGlobalState('location');
    let location = _location!;

    // so we cancel the timeout if something else changes the state
    let [formStateTimeout, setFormStateTimeout] = useState<any | null>(null);
    const submit = () => {
        setState("loading");
        scan(handle, location.id)
            .then((ev) => {
                setEvent(ev);
                setState("submitted");
                let timeout = setTimeout(() => {
                    setState("form")
                    alert("state");
                }, 4000);
                alert(timeout);
                setFormStateTimeout(() => timeout)
            })
            .catch((e: any) => {
                onCatch(e);
                setState("form");
            })
    }

    let formInputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const handleGlobalKeyPress = (event: KeyboardEvent) => {
            setState("form")
            console.log(formStateTimeout);
            if (formStateTimeout !== null) {
                clearTimeout(formStateTimeout);
                setFormStateTimeout(null);
            }

            formInputRef.current!.focus()
        }

        window.addEventListener('keydown', handleGlobalKeyPress);

        return () => window.removeEventListener('keydown', handleGlobalKeyPress)
    }, [formStateTimeout])

    const handleKeyDown = (e: any) => {
        if (e.key === "Enter" && handle) {
            submit();
        }
    }

    // If there is no location return a loading spinner
    if (!_location) {
        return <Spinner/>
    }

    // The element inside the container card
    let contentElem;
    if (state === "form" || state === "loading") {
        contentElem = <>
            <h1 className="bp3-heading">Please scan badge to sign into the {location.name}</h1>
            <div className="bp3-text-large bp3-text-muted mb-5">If you do not have a badge, contact the proctor
                on duty.
            </div>
            <FormGroup
                label="Badge ID"
                helperText="After you scan your badge, this form will submit automatically">
                <InputGroup large onChange={handleChange} onKeyDown={handleKeyDown} placeholder=""
                            id="student-handle-input" leftIcon={"align-justify"} autoComplete={"off"} spellCheck={false}
                            autoFocus inputRef={formInputRef}
                            rightElement={<Button minimal rightIcon={"arrow-right"} loading={state === "loading"}
                                                  onClick={submit}/>}/>
            </FormGroup>
        </>;
    } else if (state === "submitted") {
        contentElem = <Submitted event={event!}/>
    }

    return <div className="flex items-center justify-center m-auto"
                onFocus={() => document.getElementById("student-handle-input")!.focus()}>
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