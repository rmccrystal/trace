import React from 'react';
import {Card, ICardProps} from "@blueprintjs/core";
import {ContactReport} from "../api";
import StudentSelect from "./studentSelect";

export default function ContactTracing({...props}: ICardProps) {
    return <Card className="m-8" {...props}>
        <StudentSelect onSelect={alert} />
    </Card>
}

function ViewContactReport({report}: {report: ContactReport}) {

}