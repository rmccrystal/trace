import React, {useEffect, useState} from "react";
import {Navbar, Button} from "@blueprintjs/core";
import "./nav.scss";
import LocationSelect from "./locationSelect";
import * as Api from "../api";
import {useGlobalState} from "../app";
import {Link} from "react-router-dom";

export default function Nav({location, setLocation}: { location: Api.TraceLocation, setLocation: (location: Api.TraceLocation) => void }) {
    let [dark, setDark] = useGlobalState('dark');
    const onToggleDark = () => {
        setDark(!dark);
    }

    return <Navbar>
        <Navbar.Group align="left">
            <LocationSelect activeLocation={location} onSelect={setLocation}/>
            <Navbar.Divider/>
            <Link to="/">
                <Button minimal className="mx-1" icon="align-justify" text="Scan"/>
            </Link>
            <Link to={"/dashboard"}>
                <Button minimal className="mx-1" icon="dashboard" text={`${location.name} Dashboard`}/>
            </Link>
        </Navbar.Group>
        <Navbar.Group align="right">
            <Button icon="contrast" minimal onClick={onToggleDark}/>
        </Navbar.Group>
    </Navbar>
}