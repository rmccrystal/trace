import React, {useEffect, useState} from "react";
import {Button, Card, HTMLTable, Spinner} from "@blueprintjs/core";
import {getStudentsAtLocation, logoutStudent, Student, TraceEvent, TraceLocation} from "../api";
import {formatAMPM, onCatch, onCatchPrefix} from "./util";
import moment from "moment";
import CurrentlyInLocation from "./currentlyInLocation";

// todo: preserve state while changing the page back
export default function Dashboard({location}: {location: TraceLocation}) {
    return <CurrentlyInLocation location={location} className="max-w-xl w-full m-8" />
}
