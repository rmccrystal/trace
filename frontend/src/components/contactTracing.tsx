import React from 'react';
import {Card, ICardProps} from "@blueprintjs/core";

export default function ContactTracing({...props}: ICardProps) {
    return <Card className="m-8" {...props}>
        contact tracing
    </Card>
}