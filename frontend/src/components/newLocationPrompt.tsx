import React, {useState} from "react";
import {Button, Card, ControlGroup, FormGroup, H1, ICardProps, InputGroup} from "@blueprintjs/core";
import {createNewLocation, TraceLocation} from "../api";
import {onCatchPrefix} from "./util";

export default function NewLocationPrompt({submitCallback, ...props}: { submitCallback: (location: TraceLocation) => void } & ICardProps) {
    const [loading, setLoading] = useState(false);
    const [locationName, setLocationName] = useState("");

    const submit = () => {
        setLoading(true);
        createNewLocation({name: locationName})
            .then(submitCallback)
            .catch(onCatchPrefix("Error creating new location: "))
            .finally(() => setLoading(false));
    }

    const handleKeyDown = (e: any) => {
        if (e.key === "Enter") {
            submit();
        }
    }

    return <Card {...props}>
        <H1>No locations configured</H1>
        <FormGroup
            label="Create new location">
            <ControlGroup fill>
                <InputGroup placeholder="Location name" className="w-full" value={locationName}
                            onChange={(event: any) => setLocationName(event.target.value)} onKeyDown={handleKeyDown}/>
                <Button icon="arrow-right" loading={loading} onClick={submit}/>
            </ControlGroup>
        </FormGroup>
    </Card>
}